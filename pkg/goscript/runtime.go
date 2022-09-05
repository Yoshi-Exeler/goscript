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
		case INDEX_ASSIGN:
			r.execIndexAssign(operation)
		case CALL:
			r.execFunctionExpression(operation.Args[0].(*Expression))
		case BIND:
			r.execBind(operation)
		case RETURN:
			returnExpr := operation.Args[0].(*Expression)
			return r.unlink(r.ResolveExpression(returnExpr))
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
		case GROW:
			r.execGrow(operation)
		case SHRINK:
			r.execShrink(operation)
		default:
			panic(fmt.Sprintf("[GSR] runtime exception, invalid operation type %v", operation.Type))
		}
		r.ProgramCounter++
	}
}

func (r *Runtime) execGrow(operation *BinaryOperation) {
	// get the symbol reference from arg0
	symbolRef := operation.Args[0].(int)
	// get the amount to grow by from arg1
	amount := operation.Args[1].(int)
	// get the type from arg2
	elemType := operation.Args[2].(BinaryType)
	// fetch the array from the symbol table
	array := *r.SymbolTable[symbolRef].Value.(*[]*BinaryTypedValue)
	// save the last index of the current array
	prevLast := len(array)
	// append an empty array of *BinaryTypedValue in a single array grow
	array = append(array, make([]*BinaryTypedValue, amount)...)
	// initialize the new values correctly
	for i := prevLast; i < len(array); i++ {
		array[i] = &BinaryTypedValue{
			Type:  elemType,
			Value: defaultValuePtrOf(elemType),
		}
	}
	// save the result back to the symbol table
	r.SymbolTable[symbolRef].Value = &array
}

func (r *Runtime) execShrink(operation *BinaryOperation) {
	// get the symbol reference from arg0
	symbolRef := operation.Args[0].(int)
	// get the amount to shrink by from arg1
	amount := operation.Args[1].(int)
	// fetch the array from the symbol table
	array := *r.SymbolTable[symbolRef].Value.(*[]*BinaryTypedValue)
	// shrink the array
	newArray := array[:len(array)-amount]
	// save the result back to the symbol table
	r.SymbolTable[symbolRef].Value = &newArray
}

func (r *Runtime) unlinkedAssign(target *BinaryTypedValue, value *BinaryTypedValue) {
	switch value.Type {
	case BT_INT8:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*int8) = *value.Value.(*int8)
	case BT_INT16:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*int16) = *value.Value.(*int16)
	case BT_INT32:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*int32) = *value.Value.(*int32)
	case BT_INT64:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*int64) = *value.Value.(*int64)
	case BT_UINT8:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*uint8) = *value.Value.(*uint8)
	case BT_UINT16:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*uint16) = *value.Value.(*uint16)
	case BT_UINT32:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*uint32) = *value.Value.(*uint32)
	case BT_UINT64:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*uint64) = *value.Value.(*uint64)
	case BT_BYTE:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*byte) = *value.Value.(*byte)
	case BT_FLOAT32:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*float32) = *value.Value.(*float32)
	case BT_FLOAT64:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*float64) = *value.Value.(*float64)
	case BT_STRING:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*string) = *value.Value.(*string)
	case BT_BOOLEAN:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*bool) = *value.Value.(*bool)
	case BT_ARRAY:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*[]*BinaryTypedValue) = *value.Value.(*[]*BinaryTypedValue)
	default:
		panic("unexpected type in unlink")
	}
}

func (r *Runtime) unlink(value *BinaryTypedValue) *BinaryTypedValue {
	switch value.Type {
	case BT_INT8:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*int8)
		value.Value = &underlying
		return value
	case BT_INT16:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*int16)
		value.Value = &underlying
		return value
	case BT_INT32:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*int32)
		value.Value = &underlying
		return value
	case BT_INT64:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*int64)
		value.Value = &underlying
		return value
	case BT_UINT8:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*uint8)
		value.Value = &underlying
		return value
	case BT_UINT16:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*uint16)
		value.Value = &underlying
		return value
	case BT_UINT32:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*uint32)
		value.Value = &underlying
		return value
	case BT_UINT64:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*uint64)
		value.Value = &underlying
		return value
	case BT_BYTE:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*byte)
		value.Value = &underlying
		return value
	case BT_FLOAT32:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*float32)
		value.Value = &underlying
		return value
	case BT_FLOAT64:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*float64)
		value.Value = &underlying
		return value
	case BT_STRING:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*string)
		value.Value = &underlying
		return value
	case BT_BOOLEAN:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*bool)
		value.Value = &underlying
		return value
	case BT_ARRAY:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*[]*BinaryTypedValue)
		value.Value = &underlying
		return value
	default:
		panic("unexpected type in unlink")
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

func (r *Runtime) execBind(operation *BinaryOperation) {
	// get the symbol reference from arg0
	symbolRef := operation.Args[0].(int)
	// get the symbol type from arg0
	symType := operation.Args[1].(BinaryType)
	// initialize the symbol
	r.SymbolTable[symbolRef] = &BinaryTypedValue{
		Type:  symType,
		Value: defaultValuePtrOf(symType),
	}
	// save the symbol reference to the current scope
	r.SymbolScopeStack[len(r.SymbolScopeStack)-1] = append(r.SymbolScopeStack[len(r.SymbolScopeStack)-1], symbolRef)
}

func defaultValuePtrOf(valueType BinaryType) any {
	switch valueType {
	case BT_INT8:
		zero := int8(0)
		return &zero
	case BT_INT16:
		zero := int16(0)
		return &zero
	case BT_INT32:
		zero := int32(0)
		return &zero
	case BT_INT64:
		zero := int64(0)
		return &zero
	case BT_UINT8:
		zero := uint8(0)
		return &zero
	case BT_UINT16:
		zero := uint16(0)
		return &zero
	case BT_UINT32:
		zero := uint32(0)
		return &zero
	case BT_UINT64:
		zero := uint64(0)
		return &zero
	case BT_BYTE:
		zero := byte(0)
		return &zero
	case BT_FLOAT32:
		zero := float32(0)
		return &zero
	case BT_FLOAT64:
		zero := float64(0)
		return &zero
	case BT_STRING:
		zero := ""
		return &zero
	case BT_ARRAY:
		zero := []*BinaryTypedValue{}
		return &zero
	case BT_NOTYPE:
		return nil
	default:
		panic("invalid type")
	}
}

func (r *Runtime) execAssign(operation *BinaryOperation) {
	// get the symbol reference from arg0
	symbolRef := operation.Args[0].(int)
	// get the expression from arg1
	expression := operation.Args[1].(*Expression)
	// resolve the expression and  assign the resolution to the referenced symbol, without linking it to the expression
	r.unlinkedAssign(r.SymbolTable[symbolRef], r.ResolveExpression(expression))
}

func (r *Runtime) execIndexAssign(operation *BinaryOperation) {
	// get the symbol reference from arg0
	symbolRef := operation.Args[0].(int)
	// get the index from arg1
	index := r.ResolveExpression(operation.Args[1].(*Expression))
	// get the expression from arg2
	expression := operation.Args[2].(*Expression)
	// resolve the expression and  assign the resolution to the referenced symbol, without linking it to the expression
	r.unlinkedAssign((*r.SymbolTable[symbolRef].Value.(*[]*BinaryTypedValue))[indirectCast[int](*index)], r.ResolveExpression(expression))
}

// ResolveExpression will recursively resolve the expression to a typed value.
func (r *Runtime) ResolveExpression(e *Expression) *BinaryTypedValue {
	switch e.Operator {
	case BO_CONSTANT:
		// if the expression is constant, return its value
		return e.Value
	case BO_VSYMBOL:
		// if the expression is a variable symbol reference, just yield the symbols value
		return r.SymbolTable[e.Ref]
	case BO_FUNCTION_CALL:
		// if the expression is a function call, start executing the function until it eventually returns a constant
		return r.execFunctionExpression(e)
	case BO_BUILTIN_CALL:
		// if the expression is a function call, start executing the function until it eventually returns a constant
		return r.execFunctionExpression(e)
	case BO_INDEX_INTO:
		return r.indexIntoExpression(e)
	default:
		// otherwise, resolve the left expression
		left := r.ResolveExpression(e.LeftExpression)
		// then resolve the right expression
		right := r.ResolveExpression(e.RightExpression)
		// finally apply the operator
		return applyOperator(left, right, e.Operator, e.Value)
	}
}

// indexIntoExpression will index into the following expression, assuming it is an array and has been type checked
func (r *Runtime) indexIntoExpression(e *Expression) *BinaryTypedValue {
	// fetch the symbol from the symbol table
	symbol := r.SymbolTable[e.Ref]
	// resolve the index expression
	index := r.ResolveExpression(e.Value.Value.(*Expression))
	// index into the symbol
	return (*symbol.Value.(*[]*BinaryTypedValue))[indirectCast[int](*index)]
}

// execBuiltinCall executes the expression as a builtin function, assuming that it has been type checked before
func (r *Runtime) execBuiltinCall(e *Expression) *BinaryTypedValue {
	// convert the builtin index to a builtin index type
	builtinIdx := BuiltinFunction(e.Ref)
	// call the apropriate handler
	switch builtinIdx {
	case BF_LEN:
		return r.builtinLen(e.Args)
	case BF_INPUT:
		panic("not implemented")
	case BF_INPUTLN:
		panic("not implemented")
	case BF_MAX:
		panic("not implemented")
	case BF_MIN:
		panic("not implemented")
	case BF_PRINT:
		return r.builtinPrint(e.Args)
	case BF_PRINTF:
		panic("not implemented")
	case BF_PRINTLN:
		panic("not implemented")
	default:
		panic(fmt.Sprintf("unknown builtin %v, fatal error", builtinIdx))
	}
}

func (r *Runtime) builtinPrint(args []*FunctionArgument) *BinaryTypedValue {
	// expect the number of arguments to be 1
	expectLength(args, 1, "print builtin takes one argument")
	// perform the print
	printUnderlying(r.ResolveExpression(args[0].Expression).Value.(*BinaryTypedValue))
	// yield null
	return &BinaryTypedValue{
		Type:  BT_NOTYPE,
		Value: nil,
	}
}

// builtinLen runs the len builtin function
func (r *Runtime) builtinLen(args []*FunctionArgument) *BinaryTypedValue {
	// expect the number of arguments to be 1
	expectLength(args, 1, "length builtin takes one argument")
	// return the length of the array
	return &BinaryTypedValue{
		Value: len(*r.ResolveExpression(args[0].Expression).Value.(*[]*BinaryTypedValue)),
		Type:  BT_UINT64,
	}
}

// execFunctionExpression will execute the expression as a function, assuming that it has been type checked before
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
		r.SymbolTable[arg.SymbolRef].Value = argResolution.Value
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
