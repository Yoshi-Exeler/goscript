package goscript

import (
	"bufio"
	"fmt"
	"os"
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
		case EXPRESSION:
			r.ResolveExpression(operation.Args[0].(*Expression))
		case BIND:
			r.execBind(operation)
		case RETURN:
			if len(operation.Args) > 0 {
				returnExpr := operation.Args[0].(*Expression)
				return r.unlink(r.ResolveExpression(returnExpr))
			}
			return &BinaryTypedValue{Type: BT_NOTYPE}
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
	case BT_CHAR:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*rune) = *value.Value.(*rune)
	case BT_BOOLEAN:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*bool) = *value.Value.(*bool)
	case BT_LIST:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*[]*BinaryTypedValue) = *value.Value.(*[]*BinaryTypedValue)
	case BT_NULL:
		target.Type = BT_NULL
		target.Value = nil
	default:
		panic(fmt.Sprintf("unexpected type in unlink: %v", value.Type))
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
	case BT_CHAR:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*rune)
		value.Value = &underlying
		return value
	case BT_BOOLEAN:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*bool)
		value.Value = &underlying
		return value
	case BT_LIST:
		// cast the value's type to its underlying type
		underlying := *value.Value.(*[]*BinaryTypedValue)
		value.Value = &underlying
		return value
	case BT_NOTYPE:
		return value
	case BT_NULL:
		return value
	default:
		panic(fmt.Sprintf("unexpected type in unlink: %v", value.Type))
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
	case BT_CHAR:
		zero := rune(0)
		return &zero
	case BT_LIST:
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
	r.unlinkedAssign((*r.SymbolTable[symbolRef].Value.(*[]*BinaryTypedValue))[indirectCast[int](index)], r.ResolveExpression(expression))
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
		return r.execBuiltinCall(e)
	case BO_INDEX_INTO:
		return r.indexIntoExpression(e)
	case BO_NULLEXPR:
		return &BinaryTypedValue{
			Type: BT_NULL,
		}
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
	return (*symbol.Value.(*[]*BinaryTypedValue))[indirectCast[int](index)]
}

/*
	BF_TOUINT8   BuiltinFunction = 9
	BF_TOUINT16  BuiltinFunction = 10
	BF_TOUINT32  BuiltinFunction = 11
	BF_TOUINT64  BuiltinFunction = 12
	BF_TOINT8    BuiltinFunction = 13
	BF_TOINT16   BuiltinFunction = 14
	BF_TOINT32   BuiltinFunction = 15
	BF_TOINT64   BuiltinFunction = 16
	BF_TOFLOAT32 BuiltinFunction = 17
	BF_TOFLOAT64 BuiltinFunction = 18
	BF_TOSTRING  BuiltinFunction = 19
	BF_TOCHAR    BuiltinFunction = 20
	BF_TOBYTE    BuiltinFunction = 21
*/

func (r *Runtime) builtinToUint8(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToUint8 takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[uint8](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_UINT8,
		Value: &conv,
	}
}

func (r *Runtime) builtinToUint16(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToUint16 takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[uint16](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_UINT16,
		Value: &conv,
	}
}

func (r *Runtime) builtinToUint32(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToUint32 takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[uint32](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_UINT32,
		Value: &conv,
	}
}

func (r *Runtime) builtinToUint64(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToUint64 takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[uint64](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_UINT64,
		Value: &conv,
	}
}

func (r *Runtime) builtinToint8(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToInt8 takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[int8](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_INT8,
		Value: &conv,
	}
}

func (r *Runtime) builtinToint16(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToInt16 takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[int16](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_INT16,
		Value: &conv,
	}
}

func (r *Runtime) builtinToint32(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToInt32 takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[int32](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_INT32,
		Value: &conv,
	}
}

func (r *Runtime) builtinToint64(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToInt64 takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[int64](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_INT64,
		Value: &conv,
	}
}

func (r *Runtime) builtinTofloat32(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToFloat32 takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[float32](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_FLOAT32,
		Value: &conv,
	}
}

func (r *Runtime) builtinTofloat64(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToFloat64 takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[float64](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_FLOAT64,
		Value: &conv,
	}
}

func (r *Runtime) builtinToByte(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToByte takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[byte](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_BYTE,
		Value: &conv,
	}
}

func (r *Runtime) builtinToString(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToString takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := sprintUnderlying(r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_STRING,
		Value: &conv,
	}
}

func (r *Runtime) builtinToChar(args []*FunctionArgument) *BinaryTypedValue {
	// expect one argument
	expectLength(args, 1, "builtinToChar takes one argument")
	// resolve the expression to a value and cast to uint8
	conv := indirectCast[rune](r.ResolveExpression(args[0].Expression))
	return &BinaryTypedValue{
		Type:  BT_CHAR,
		Value: &conv,
	}
}

// execBuiltinCall executes the expression as a builtin function, assuming that it has been type checked before
func (r *Runtime) execBuiltinCall(e *Expression) *BinaryTypedValue {
	// convert the builtin index to a builtin index type
	builtinIdx := BuiltinFunction(e.Ref)
	// call the appropriate handler
	switch builtinIdx {
	case BF_LEN:
		return r.builtinLen(e.Args)
	case BF_INPUT:
		return r.builtinInput()
	case BF_INPUTLN:
		return r.builtinInputln()
	case BF_MAX:
		panic("not implemented")
	case BF_MIN:
		panic("not implemented")
	case BF_PRINT:
		return r.builtinPrint(e.Args)
	case BF_PRINTF:
		return r.builtinPrintf(e.Args)
	case BF_PRINTLN:
		return r.builtinPrintln(e.Args)
	case BF_TOUINT8:
		return r.builtinToUint8(e.Args)
	case BF_TOUINT16:
		return r.builtinToUint16(e.Args)
	case BF_TOUINT32:
		return r.builtinToUint32(e.Args)
	case BF_TOUINT64:
		return r.builtinToUint64(e.Args)
	case BF_TOINT8:
		return r.builtinToint8(e.Args)
	case BF_TOINT16:
		return r.builtinToint16(e.Args)
	case BF_TOINT32:
		return r.builtinToint32(e.Args)
	case BF_TOINT64:
		return r.builtinToint64(e.Args)
	case BF_TOFLOAT32:
		return r.builtinTofloat32(e.Args)
	case BF_TOFLOAT64:
		return r.builtinTofloat64(e.Args)
	case BF_TOBYTE:
		return r.builtinToByte(e.Args)
	case BF_TOSTRING:
		return r.builtinToString(e.Args)
	case BF_TOCHAR:
		return r.builtinToChar(e.Args)
	default:
		panic(fmt.Sprintf("unknown builtin %v, fatal error", builtinIdx))
	}
}

func (r *Runtime) builtinInput() *BinaryTypedValue {
	buff := make([]byte, 1)
	os.Stdin.Read(buff)
	val := rune(buff[0])
	return &BinaryTypedValue{
		Type:  BT_CHAR,
		Value: &val,
	}
}

func (r *Runtime) builtinInputln() *BinaryTypedValue {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Sprintf("cannot read from stdin with error %v", err))
	}
	return &BinaryTypedValue{
		Type:  BT_STRING,
		Value: &text,
	}
}

func (r *Runtime) builtinPrint(args []*FunctionArgument) *BinaryTypedValue {
	// expect the number of arguments to be 1
	expectLength(args, 1, "print builtin takes one argument")
	// perform the print
	printUnderlying(r.ResolveExpression(args[0].Expression))
	// yield null
	return &BinaryTypedValue{
		Type:  BT_NOTYPE,
		Value: nil,
	}
}

func (r *Runtime) builtinPrintln(args []*FunctionArgument) *BinaryTypedValue {
	// expect the number of arguments to be 1
	expectLength(args, 1, "println builtin takes one argument")
	// perform the print
	printlnUnderlying(r.ResolveExpression(args[0].Expression))
	// yield null
	return &BinaryTypedValue{
		Type:  BT_NOTYPE,
		Value: nil,
	}
}

func (r *Runtime) builtinPrintf(args []*FunctionArgument) *BinaryTypedValue {
	// resolve the other args
	argSpread := r.resolveArgs(args[1:])
	// perform the print
	printfUnderlying(*(r.ResolveExpression(args[0].Expression).Value.(*string)), argSpread)
	// yield null
	return &BinaryTypedValue{
		Type:  BT_NOTYPE,
		Value: nil,
	}
}

func (r *Runtime) resolveArgs(args []*FunctionArgument) []any {
	res := []any{}
	for _, arg := range args {
		res = append(res, dereferenceUnderlying(r.ResolveExpression(arg.Expression)))
	}
	return res
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
