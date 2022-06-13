package goscript

import (
	"fmt"
)

func NewRuntime() *Runtime {
	return &Runtime{}
}

type Runtime struct {
	SymbolTable      []*BinaryTypedValue
	SymbolScopeStack [][]int // [scope depth][symbols]
	ProgramCounter   int
	Program          Program
}

// reset will reset the state of the runtime
func (r *Runtime) reset() {
	r.ProgramCounter = 0
	r.SymbolTable = []*BinaryTypedValue{}
	r.SymbolScopeStack = make([][]int, 1)
}

func (r *Runtime) enterScope() {
	// place a new scope on the scope stack
	r.SymbolScopeStack = append(r.SymbolScopeStack, []int{})
}

// exitScope will cleanup symbols when leaving a scope
func (r *Runtime) exitScope() {
	// grab the symbols references from the top of the symbol scope stack
	symbolsToFree := r.SymbolScopeStack[len(r.SymbolScopeStack)-1]
	// free the symbols referenced here
	for _, symbolRef := range symbolsToFree {
		// just nil the elements so go can gc them for us
		r.SymbolTable[symbolRef] = nil
	}
	// pop the scope stack
	r.SymbolScopeStack = r.SymbolScopeStack[:len(r.SymbolScopeStack)-1]
}

// Exec will reset the runtime and then run the specified program until it completes
func (r *Runtime) Exec(program Program) any {
	// completely reset the runtime
	r.reset()
	// save our program
	r.Program = program
	// build a symbol table of the requested size
	r.SymbolTable = make([]*BinaryTypedValue, program.SymbolTableSize)
	// execute our program until the main function returns
	returnValue := r.execUntilReturn()
	// exit with this value
	return returnValue
}

// execUntilReturn will keep executing instructions until a return is hit in the current scope, and then return the value passed to the return
func (r *Runtime) execUntilReturn() *BinaryTypedValue {
	for {
		// fetch the next operation from the program
		operation := &r.Program.Operations[r.ProgramCounter]
		// fmt.Printf("[%v] %v\n", r.ProgramCounter, operation.String())
		// time.Sleep(time.Millisecond * 50)
		// call the operation handler for this operation type
		switch operation.Type {
		case ASSIGN:
			r.execAssign(operation)
		case CALL:
			r.execCall(operation)
		case BIND:
			r.execBind(operation)
		case RETURN:
			returnExpr := operation.Args[0].(*Expression)
			return r.ResolveExpression(returnExpr)
		case ENTER_SCOPE:
			r.enterScope()
		case EXIT_SCOPE:
			r.exitScope()
		case JUMP:
			r.execJump(operation)
		case JUMP_IF:
			r.execJumpIf(operation)
		case JUMP_IF_NOT:
			r.execJumpIfNot(operation)
		default:
			panic(fmt.Sprintf("[GSR] runtime exception, invalid operation type %v", operation.Type))
		}
		r.ProgramCounter++
	}
}

func (r *Runtime) execJumpIf(operation *BinaryOperation) {
	// get the expression from arg0
	condition := operation.Args[0].(*Expression)
	// resolve the expressions
	resolution := r.ResolveExpression(condition).Value.(bool)
	// if the condition is true, we jump
	if resolution {
		// get the target address from arg1
		targetPC := operation.Args[1].(int)
		// jump to the target
		r.ProgramCounter = targetPC
	}
}

func (r *Runtime) execJumpIfNot(operation *BinaryOperation) {
	// get the expression from arg0
	condition := operation.Args[0].(*Expression)
	// resolve the expressions
	resolution := r.ResolveExpression(condition).Value.(bool)
	// if the condition is true, we jump
	if !resolution {
		// get the target address from arg1
		targetPC := operation.Args[1].(int)
		// jump to the target
		r.ProgramCounter = targetPC
	}
}

func (r *Runtime) execJump(operation *BinaryOperation) {
	// get the target address from arg0
	targetPC := operation.Args[0].(int)
	// jump to the target
	r.ProgramCounter = targetPC
}

// execCall is a wrapper that runs a function call
func (r *Runtime) execCall(operation *BinaryOperation) {
	// get the function expression from arg0
	_ = r.execFunctionExpression(operation.Args[0].(*Expression))
}

func (r *Runtime) execBind(operation *BinaryOperation) {
	// get the symbol reference from arg0
	symbolRef := operation.Args[0].(int)
	// initialize the symbol
	r.SymbolTable[symbolRef] = &BinaryTypedValue{}
	// save the symbol reference to the current scope
	r.SymbolScopeStack[len(r.SymbolScopeStack)-1] = append(r.SymbolScopeStack[len(r.SymbolScopeStack)-1], symbolRef)
}

func (r *Runtime) execAssign(operation *BinaryOperation) {
	// get the symbol reference from arg0
	symbolRef := operation.Args[0].(int)
	// get the expression from arg1
	expression := operation.Args[1].(*Expression)
	// resolve the expression and  assign the resolution to the referenced symbol
	resolution := r.ResolveExpression(expression)
	sym := r.SymbolTable[symbolRef]
	sym.Value = resolution.Value
	sym.Type = resolution.Type
}

// ResolveExpression will recursively resolve the expression to a typed value.
func (r *Runtime) ResolveExpression(e *Expression) *BinaryTypedValue {
	// if the expression is constant, return its value
	if e.IsConstant() {
		return e.Value
	}
	// if the expression is a variable symbol reference, just yield the symbols value
	if e.isVSymbol() {
		return r.SymbolTable[e.Ref]
	}
	// if the expression is a function call, start executing the function until it eventually returns a constant
	if e.IsFunction() {
		return r.execFunctionExpression(e)
	}
	// otherwise, resolve the left expression
	left := r.ResolveExpression(e.LeftExpression)
	// then resolve the right expression
	right := r.ResolveExpression(e.RightExpression)
	// finally apply the operator
	return applyOperator(left, right, e.Operator, e.Value)
}

// exec will execute the expression as a function, assuming that it has been type checked before
func (r *Runtime) execFunctionExpression(e *Expression) *BinaryTypedValue {
	// save the current pc so we can return here later
	returnPC := r.ProgramCounter
	// jump to the appropriate section
	r.ProgramCounter = e.Ref
	// open a new scope
	r.enterScope()
	// perform the appropriate argument mapping
	for _, arg := range e.Args {
		// resolve the argument expression
		argResolution := r.ResolveExpression(arg.Expression)
		// set the symbol in the local scope
		r.SymbolTable[arg.SymbolRef].Value = argResolution
	}
	// execute until this top level function returns
	e.Value = r.execUntilReturn()
	// exit the scope
	r.exitScope()
	// return to the original place in the code
	r.ProgramCounter = returnPC
	return e.Value
}

func (r *Runtime) defaultValueOf(binaryType BinaryType) any {
	switch binaryType {
	case BT_INT8:
		return int8(0)
	case BT_INT16:
		return int16(0)
	case BT_INT32:
		return int32(0)
	case BT_INT64:
		return int64(0)
	case BT_UINT8:
		return uint8(0)
	case BT_UINT16:
		return uint16(0)
	case BT_UINT32:
		return uint32(0)
	case BT_UINT64:
		return uint64(0)
	case BT_STRING:
		return ""
	case BT_CHAR:
		return ""
	case BT_BYTE:
		return byte(0)
	case BT_FLOAT32:
		return float32(0)
	case BT_FLOAT64:
		return float64(0)
	case BT_ANY:
		return nil
	case BT_STRUCT:
		return nil
	default:
		panic(fmt.Sprintf("[GSR] runtime exception, invalid symbol type %v", binaryType))
	}
}
