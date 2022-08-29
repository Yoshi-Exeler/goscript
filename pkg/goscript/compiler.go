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
  - Resolve constant expressions as far as possible

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
		if !c.calledFunctionByName[name] && name != "main" {
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
	fmt.Printf("[GSC][STAGE_COMPLETION] generating bytecode completed in %v\n", time.Since(startGenBytecode))
	fmt.Printf("[GSC][generateProgram] completed in %v\n", time.Since(start))
	return nil, nil
}

func (c *Compiler) compileFunction(def *FunctionDefinition) []Operation {
	//
	return nil
}

func (c *Compiler) prescanFunction(def *FunctionDefinition) {
	// scan the parameters
	for _, param := range def.Accepts {
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

func (c *Compiler) scanExpression(expr *Expression) {
	if expr.Operator == BO_FUNCTION_CALL_PLACEHOLDER {
		wasKnown := c.calledFunctionByName[expr.Value.Value.(*FunctionCallPlaceholder).Name]
		c.calledFunctionByName[expr.Value.Value.(*FunctionCallPlaceholder).Name] = true
		if !wasKnown {
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
