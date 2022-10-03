package goscript

import (
	"fmt"
	"time"
)

type Compiler struct {
	funcsByName          map[string]*FunctionDefinition
	funcBaseByName       map[string]int
	symbolByName         map[string]*IntermediateVar
	symbolIndexByName    map[string]int
	currentSymbolIndex   int
	calledFunctionByName map[string]bool
	currentProgram       *Program
	currentOpIndex       int
	currentFunction      *FunctionDefinition
}

func NewCompiler() *Compiler {
	return &Compiler{
		funcsByName:          make(map[string]*FunctionDefinition),
		funcBaseByName:       make(map[string]int),
		symbolByName:         make(map[string]*IntermediateVar),
		symbolIndexByName:    make(map[string]int),
		calledFunctionByName: make(map[string]bool),
		currentSymbolIndex:   0,
		currentProgram: &Program{
			Operations: []BinaryOperation{},
		},
	}
}

type CompileJob struct {
	MainFilePath       string
	VendorPath         string
	LocalWorkspaceRoot string
	StandardLibPath    string
}

func (c *Compiler) Compile(job CompileJob) (*Program, error) {
	appSource, err := discoverSources(job.MainFilePath, job.LocalWorkspaceRoot)
	if err != nil {
		return nil, err
	}
	fqsc, err := generateFQSC(appSource)
	if err != nil {
		return nil, err
	}
	intermediate := parse(fqsc)
	return c.generateProgram(intermediate)
}

/*
generateProgram generates a program from the intermediary program representation
The following steps will be performed:
- Replace Symbol Placeholders in expressions
- Replace Function Placeholders in expressions
- Eliminate Dead code
- Optimize:
  - Resolve constant expressions as far as possible (WIP)

- Generate the actual bytecode
*/
func (c *Compiler) generateProgram(intermediate *IntermediateProgram) (*Program, error) {
	fmt.Println("[GSC][generateProgram] begin generating program")
	start := time.Now()
	// map our all the functions by name
	for _, funcDef := range intermediate.Functions {
		funcDef := funcDef
		c.funcsByName[funcDef.Name] = funcDef
	}
	// start off with a prescan of the program, discovering all symbols and function calls
	fmt.Println("[GSC][codePreScan] begin code prescan")
	startScan := time.Now()
	c.prescanFunction(&intermediate.Entrypoint)
	fmt.Printf("[GSC][STAGE_COMPLETION] code prescan completed in %v\n", time.Since(startScan))
	// eliminate functions that are never called
	fmt.Println("[GSC][DCE] begin dead code elimination")
	startDce := time.Now()
	newFuncs := make(map[string]*FunctionDefinition)
	for name, called := range c.funcsByName {
		called := called
		if !c.calledFunctionByName[name] && name != "#fn_0_main_main" {
			fmt.Printf("[GSC][DCE] eliminate function %v\n", name)
			continue
		}
		newFuncs[name] = called
	}
	fmt.Printf("[GSC][STAGE_COMPLETION] DCE completed in %v\n", time.Since(startDce))
	// generate our bytecode
	fmt.Println("[GSC][generateBytecode] begin generating bytecode")
	startGenBytecode := time.Now()
	c.compileFunction(&intermediate.Entrypoint)
	for _, function := range newFuncs {
		function := function
		c.compileFunction(function)
	}
	fmt.Printf("[GSC][STAGE_COMPLETION] generating bytecode completed in %v\n", time.Since(startGenBytecode))
	fmt.Println("[GSC][finalizeProgram] begin finalizing references")
	startFinalize := time.Now()
	c.currentProgram = c.finalizeProgram(c.currentProgram)
	fmt.Printf("[GSC][STAGE_COMPLETION] finalization completed in %v\n", time.Since(startFinalize))
	fmt.Printf("[GSC][generateProgram] completed in %v\n", time.Since(start))
	c.currentProgram.SymbolTableSize = len(c.symbolIndexByName) + 1
	return c.currentProgram, nil
}

// finalizeProgram will recompile all expressions after all functions have been initially generated to fix
// all function base address references
func (c *Compiler) finalizeProgram(prog *Program) *Program {
	for i := 0; i < len(prog.Operations); i++ {
		switch prog.Operations[i].Type {
		case ASSIGN:
			prog.Operations[i] = c.recompileArgs(prog.Operations[i], 1)
		case INDEX_ASSIGN:
			prog.Operations[i] = c.recompileArgs(prog.Operations[i], 1, 2)
		case EXPRESSION:
			prog.Operations[i] = c.recompileArgs(prog.Operations[i], 0)
		case BIND:
			continue
		case RETURN:
			prog.Operations[i] = c.recompileArgs(prog.Operations[i], 0)
		case ENTER_SCOPE:
			continue
		case EXIT_SCOPE:
			continue
		case JUMP:
			continue
		case JUMP_IF:
			prog.Operations[i] = c.recompileArgs(prog.Operations[i], 0)
		case JUMP_IF_NOT:
			prog.Operations[i] = c.recompileArgs(prog.Operations[i], 0)
		case GROW:
			// TODO: grow should grow based on an expression amount
			continue
		case SHRINK:
			// TODO: shrink should shrink based on an expression amount
			continue
		default:
			panic(fmt.Sprintf("[GSC] encountered invalid expression in finalize %v", prog.Operations[i]))
		}
	}
	return prog
}

func (c *Compiler) recompileArgs(op BinaryOperation, indices ...int) BinaryOperation {
	for _, idx := range indices {
		op.Args[idx] = c.compileExpression(op.Args[idx].(*Expression))
	}
	return op
}

func (c *Compiler) compileFunction(def *FunctionDefinition) {
	c.currentFunction = def
	c.currentOpIndex = 0
	c.funcBaseByName[def.Name] = len(c.currentProgram.Operations)
	for {
		if len(c.currentFunction.Operations)-1 == c.currentOpIndex {
			break
		}
		op := c.currentFunction.Operations[c.currentOpIndex]
		switch op.Type {
		case IM_ASSIGN:
			c.generateAssign(op)
		case IM_BREAK:
			panic("break is not implemented in compileFunction")
		case IM_CLOSING_BRACKET:
			panic("should not reach closig bracket outside of inner parser")
		case IM_EXPRESSION:
			c.generateExpression(op)
		case IM_FOR:
			c.generateLoop(op)
		case IM_FOREACH:
			panic("foreach is not implemented in compileFunction")
		case IM_RETURN:
			c.generateReturn(op)
		case IM_NOP:
		default:
			panic(fmt.Sprintf("unkndown operation %v cannot compile", op.Type))
		}
		c.currentOpIndex++
	}
	// if the function does not end in a return, we will insert one
	if c.currentProgram.Operations[len(c.currentProgram.Operations)-1].Type != RETURN {
		c.currentProgram.Operations = append(c.currentProgram.Operations, BinaryOperation{
			Type: RETURN,
			Args: []any{&Expression{Operator: BO_NULLEXPR}},
		})
	}
}

func (c *Compiler) generateLoop(op *IntermediateOperation) {
	iteratorRef := c.symbolIndexByName[op.Args[0].(string)]
	c.currentProgram.Operations = append(c.currentProgram.Operations, NewEnterScope())
	c.currentProgram.Operations = append(c.currentProgram.Operations, NewBindOp(iteratorRef, op.Args[1].(IntermediateType).Type))
	c.currentProgram.Operations = append(c.currentProgram.Operations, NewAssignExpressionOp(iteratorRef, c.compileExpression(op.Args[2].(*Expression))))
	c.currentProgram.Operations = append(c.currentProgram.Operations, NewJumpIfNotOp(1, c.compileExpression(op.Args[3].(*Expression))))
	loopHeadAddr := len(c.currentProgram.Operations) - 1
	c.currentOpIndex++
	c.generateUntilClose()
	c.currentProgram.Operations = append(c.currentProgram.Operations, NewJumpOp(loopHeadAddr))
	c.currentProgram.Operations = append(c.currentProgram.Operations, NewExitScopeOp())
	loopEndAddr := len(c.currentProgram.Operations) - 1
	c.currentProgram.Operations[loopHeadAddr] = NewJumpIfNotOp(loopEndAddr, c.compileExpression(op.Args[3].(*Expression)))
}

func (c *Compiler) compileExpression(expr *Expression) *Expression {
	return c.resolveSymbols(c.resolveCalls(expr))
}

func (c *Compiler) resolveSymbols(expr *Expression) *Expression {
	if expr.Operator == BO_VSYMBOL_PLACEHOLDER {
		expr.Operator = BO_VSYMBOL
		expr.Ref = c.symbolIndexByName[expr.Value.Value.(string)]
		expr.Value = nil
	}
	for i := 0; i < len(expr.Args); i++ {
		expr.Args[i].Expression = c.resolveSymbols(expr.Args[i].Expression)
	}
	if expr.LeftExpression != nil {
		expr.LeftExpression = c.resolveSymbols(expr.LeftExpression)
	}
	if expr.RightExpression != nil {
		expr.RightExpression = c.resolveSymbols(expr.RightExpression)
	}
	return expr
}

func (c *Compiler) resolveCalls(expr *Expression) *Expression {
	if expr.Operator == BO_FUNCTION_CALL_PLACEHOLDER {
		expr.Operator = BO_FUNCTION_CALL
		// check if we have a base address for this function
		if c.funcBaseByName[expr.Value.Value.(*FunctionCallPlaceholder).Name] == 0 {
			// check if this is a builtin
			if builtins[expr.Value.Value.(*FunctionCallPlaceholder).Name] != 0 {
				expr.Operator = BO_BUILTIN_CALL
				expr.Ref = int(builtins[expr.Value.Value.(*FunctionCallPlaceholder).Name])
				// map the function arguments into the bytecode format, resolving vsymbol references
				ph := expr.Value.Value.(*FunctionCallPlaceholder)
				funcArgs := []*FunctionArgument{}
				for i := 0; i < len(ph.Args); i++ {
					funcArgs = append(funcArgs, &FunctionArgument{
						Expression: c.compileExpression(ph.Args[i]),
						SymbolRef:  c.symbolIndexByName[ph.Name],
					})
				}
				expr.Args = funcArgs
			} else {
				// if no base address exists for this function, check our known functions
				sideFunc := c.funcsByName[expr.Value.Value.(*FunctionCallPlaceholder).Name]
				if sideFunc == nil {
					panic(fmt.Sprintf("cannot call undefined function %v", expr.Value.Value.(*FunctionCallPlaceholder).Name))
				}
				// otherwise, this function is known, but not yet compiled, in which case we just exit here
				// for now and fix this in a later pass. This is significantly easier that generating a resolution
				// order graph, which also might not always work (cyclic calls, recursion)
				// in this case we also set this call back to placeholder status, since we want later passes to fix it
				expr.Operator = BO_FUNCTION_CALL_PLACEHOLDER
			}
		} else {
			// if this function is still a placeholder but we know an address for it, insert it now
			expr.Ref = c.funcBaseByName[expr.Value.Value.(*FunctionCallPlaceholder).Name]
			// map the function arguments into the bytecode format, resolving vsymbol references
			ph := expr.Value.Value.(*FunctionCallPlaceholder)
			funcArgs := []*FunctionArgument{}
			for i := 0; i < len(ph.Args); i++ {
				funcArgs = append(funcArgs, &FunctionArgument{
					Expression: c.compileExpression(ph.Args[i]),
					SymbolRef:  c.symbolIndexByName[ph.Name],
				})
			}
			expr.Args = funcArgs
		}
	}
	for i := 0; i < len(expr.Args); i++ {
		expr.Args[i].Expression = c.resolveCalls(expr.Args[i].Expression)
	}
	if expr.LeftExpression != nil {
		expr.LeftExpression = c.resolveCalls(expr.LeftExpression)
	}
	if expr.RightExpression != nil {
		expr.RightExpression = c.resolveCalls(expr.RightExpression)
	}
	return expr
}

func (c *Compiler) generateUntilClose() {
	for {
		op := c.currentFunction.Operations[c.currentOpIndex]
		switch op.Type {
		case IM_ASSIGN:
			c.generateAssign(op)
		case IM_BREAK:
			panic("break is not implemented in generateUntilClose")
		case IM_CLOSING_BRACKET:
			return
		case IM_EXPRESSION:
			c.generateExpression(op)
		case IM_FOR:
			c.generateLoop(op)
		case IM_FOREACH:
			panic("foreach is not implemented in generateUntilClose")
		case IM_RETURN:
			c.generateReturn(op)
		case IM_NOP:
		default:
			panic(fmt.Sprintf("unkndown operation %v cannot compile", op.Type))
		}
		c.currentOpIndex++
	}
}

func (c *Compiler) generateReturn(op *IntermediateOperation) {
	c.currentProgram.Operations = append(c.currentProgram.Operations, NewReturnValueOp(c.compileExpression(op.Args[0].(*Expression))))
}

func (c *Compiler) generateExpression(op *IntermediateOperation) {
	c.currentProgram.Operations = append(c.currentProgram.Operations, NewExpressionOp(c.compileExpression(op.Args[0].(*Expression))))
}

func (c *Compiler) generateAssign(op *IntermediateOperation) {
	c.currentProgram.Operations = append(c.currentProgram.Operations, NewBindOp(c.symbolIndexByName[op.Args[0].(string)], op.Args[1].(IntermediateType).Type))
	if len(op.Args) == 3 {
		c.currentProgram.Operations = append(c.currentProgram.Operations, NewAssignExpressionOp(c.symbolIndexByName[op.Args[0].(string)], c.compileExpression(op.Args[2].(*Expression))))
	} else {
		c.currentProgram.Operations = append(c.currentProgram.Operations, NewAssignExpressionOp(c.symbolIndexByName[op.Args[0].(string)], NewConstantExpression(defaultValuePtrOf(op.Args[1].(IntermediateType).Type), op.Args[1].(IntermediateType).Type)))
	}
}

func (c *Compiler) prescanFunction(def *FunctionDefinition) {
	// scan the parameters
	for _, param := range def.Accepts {
		param := param
		c.symbolByName[param.Name] = &param
		c.symbolIndexByName[param.Name] = c.currentSymbolIndex
		c.currentSymbolIndex++
	}
	// scan the operations for any symbols and function calls
	for _, op := range def.Operations {
		switch op.Type {
		case IM_ASSIGN:
			c.symbolByName[op.Args[0].(string)] = &IntermediateVar{
				Name: op.Args[0].(string),
				Type: op.Args[1].(IntermediateType),
			}
			c.symbolIndexByName[op.Args[0].(string)] = c.currentSymbolIndex
			c.currentSymbolIndex++
			// check if an assigned expression is present
			if len(op.Args) == 3 {
				c.scanExpression(op.Args[2].(*Expression))
			}
		case IM_FOR:
			c.symbolByName[op.Args[0].(string)] = &IntermediateVar{
				Name: op.Args[0].(string),
				Type: op.Args[1].(IntermediateType),
			}
			c.symbolIndexByName[op.Args[0].(string)] = c.currentSymbolIndex
			c.currentSymbolIndex++
			// scan both expressions of the loop
			c.scanExpression(op.Args[2].(*Expression))
			c.scanExpression(op.Args[3].(*Expression))
		// case IM_FOREACH:
		// 	c.symbolByName[op.Args[0].(string)] = &IntermediateVar{
		// 		Name: op.Args[0].(string),
		// 		Type: op.Args[1].(IntermediateType),
		// 	}
		// 	c.symbolIndexByName[op.Args[0].(string)] = c.currentSymbolIndex
		// 	c.currentSymbolIndex++
		// we need type infenerce here
		case IM_RETURN:
			conv, ok := op.Args[0].(*Expression)
			if ok {
				c.scanExpression(conv)
			}
		case IM_EXPRESSION:
			c.scanExpression(op.Args[0].(*Expression))
		}
	}
}

var builtins = map[string]BuiltinFunction{
	"input":   BF_INPUT,
	"inputln": BF_INPUTLN,
	"len":     BF_LEN,
	"max":     BF_MAX,
	"min":     BF_MIN,
	"print":   BF_PRINT,
	"printf":  BF_PRINTF,
	"println": BF_PRINTLN,
	"byte":    BF_TOBYTE,
	"i8":      BF_TOINT8,
	"i16":     BF_TOINT16,
	"i32":     BF_TOINT32,
	"i64":     BF_TOINT64,
	"u8":      BF_TOUINT8,
	"u16":     BF_TOUINT16,
	"u32":     BF_TOUINT32,
	"u64":     BF_TOUINT64,
	"f32":     BF_TOFLOAT32,
	"f64":     BF_TOFLOAT64,
	"char":    BF_TOCHAR,
	"str":     BF_TOSTRING,
}

func (c *Compiler) scanExpression(expr *Expression) {
	if expr.Operator == BO_FUNCTION_CALL_PLACEHOLDER {
		wasKnown := c.calledFunctionByName[expr.Value.Value.(*FunctionCallPlaceholder).Name]
		c.calledFunctionByName[expr.Value.Value.(*FunctionCallPlaceholder).Name] = true
		if !wasKnown && builtins[expr.Value.Value.(*FunctionCallPlaceholder).Name] == 0 {
			c.prescanFunction(c.funcsByName[expr.Value.Value.(*FunctionCallPlaceholder).Name])
		}
	}
	if expr.LeftExpression != nil {
		c.scanExpression(expr.LeftExpression)
	}
	if expr.RightExpression != nil {
		c.scanExpression(expr.RightExpression)
	}
}
