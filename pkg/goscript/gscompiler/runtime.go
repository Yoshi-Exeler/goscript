package gscompiler

import (
	"fmt"
	"os"
	"sync"
)

var once sync.Once

var instance *Runtime

func Init() {
	once.Do(func() {
		// do init here later
		instance = &Runtime{}
	})
}

func GetInstance() *Runtime {
	return instance
}

type Runtime struct {
	SymbolTable      []*BinarySymbol
	SymbolScopeStack [][]int // [scope depth][symbols]
	ProgramCounter   int
	Program          Program
}

// reset will reset the state of the runtime
func (r *Runtime) reset() {
	r.ProgramCounter = 0
	r.SymbolTable = []*BinarySymbol{}
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
func (r *Runtime) Exec(program Program) {
	// completely reset the runtime
	r.reset()
	// save our program
	r.Program = program
	// execute our program until the main function returns
	returnValue := r.execUntilReturn()
	// exit with this value
	r.exitWithValue(returnValue)
}

func (r *Runtime) exitWithValue(v any) {
	fmt.Printf("Program exited with %v", v)
	os.Exit(0)
}

// execUntilReturn will keep executing instructions until a return is hit in the current scope, and then return the value passed to the return
func (r *Runtime) execUntilReturn() *BinaryTypedValue {
	for {
		// fetch the next operation from the program
		operation := &r.Program[r.ProgramCounter]
		// call the operation handler for this operation type
		switch operation.Type {
		case ASSIGN_EXPRESSION:
			r.execAssignExpression(operation)
		case CONDITIONAL_BLOCK_ENTER:
			r.execConditionalBlockEnter(operation)
		case CALL_FUNCTION:
			r.execFunctionOperation(operation)
		case RETURN_VALUE:
			return operation.Args[0].(*BinaryTypedValue)
		default:
			panic(fmt.Sprintf("[GSR] runtime exception, invalid operation type %v", operation.Type))
		}
		r.ProgramCounter++
	}
}

// execUntilReturn will keep executing instructions until a return is hit in the current scope, and then return the value passed to the return
func (r *Runtime) execUntilScopeClose() {
	for {
		// fetch the next operation from the program
		operation := &r.Program[r.ProgramCounter]
		// call the operation handler for this operation type
		switch operation.Type {
		case ASSIGN_EXPRESSION:
			r.execAssignExpression(operation)
		case CONDITIONAL_BLOCK_ENTER:
			r.execConditionalBlockEnter(operation)
		case CALL_FUNCTION:
			r.execFunctionOperation(operation)
		case CLOSE_SCOPE:
			return
		default:
			panic(fmt.Sprintf("[GSR] runtime exception, invalid operation type %v", operation.Type))
		}
		r.ProgramCounter++
	}
}

// execFunctionOperation is a wrapper that runs a function call
func (r *Runtime) execFunctionOperation(operation *BinaryOperation) {
	// get the function expression from arg0
	r.execFunctionExpression(operation.Args[0].(*Expression))
}

func (r *Runtime) execAssignExpression(operation *BinaryOperation) {
	// get the symbol reference from arg0
	symbolRef := operation.Args[0].(int)
	// get the expression from arg1
	expression := operation.Args[1].(*Expression)
	// resolve the expression
	resolution := r.ResolveExpression(expression)
	// assign the resolution to the referenced symbol
	r.SymbolTable[symbolRef].Value = resolution
}

func (r *Runtime) execConditionalBlockEnter(operation *BinaryOperation) {
	// get the conditional expression from arg0
	condition := operation.Args[0].(*Expression)
	// resolve the expression
	resolution := r.ResolveExpression(condition).Value.(bool)
	// if the condition resolved to true, we should enter the if branch
	if resolution {
		// get the main branch entry pc from arg1
		blockPC := operation.Args[1].(int)
		// jump to the block pc
		r.ProgramCounter = blockPC
		// enter a new scope
		r.enterScope()
		// execute until this scope completes
		r.execUntilScopeClose()
		// exit the scope
		r.exitScope()
		// get the post condition pc from arg2
		nextPC := operation.Args[2].(int)
		// jump to the next 'real' instruction
		r.ProgramCounter = nextPC
		return
	}
	// check if there is an else block pc in arg3
	elsePC, ok := operation.Args[3].(int)
	if !ok || elsePC == 0 {
		// if there is no else block pc we just jump to the next real address and exit
		// get the post condition pc from arg2
		nextPC := operation.Args[2].(int)
		// jump to the next 'real' instruction
		r.ProgramCounter = nextPC
		return
	}
	// if there is a valid else block pc, jump into it
	r.ProgramCounter = elsePC
	// enter a new scope
	r.enterScope()
	// execute until this scope completes
	r.execUntilScopeClose()
	// exit the scope
	r.exitScope()
	// get the post condition pc from arg2
	nextPC := operation.Args[2].(int)
	// jump to the next 'real' instruction
	r.ProgramCounter = nextPC
}

func (r *Runtime) ResolveExpression(e *Expression) *BinaryTypedValue {
	// if the expression is constant, return its value
	if e.IsConstant() {
		return &BinaryTypedValue{
			Type:  e.Type,
			Value: e.Value,
		}
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
	return applyOperator(left, right, e.Operator)
}

// exec will execute the expression as a function, assuming that it has been type checked before
func (r *Runtime) execFunctionExpression(e *Expression) *BinaryTypedValue {
	// fetch the call information
	call := e.Value.(*BinaryFunctionCall)
	// save the current pc so we can return here later
	returnPC := r.ProgramCounter
	// jump to the appropriate section
	r.ProgramCounter = call.BlockEntry
	// open a new scope
	r.enterScope()
	// perform the appropriate argument mapping
	for outer, inner := range call.Args {
		// copy the symbol to avoid mutating the external state
		localSymbol := *r.SymbolTable[outer]
		// set the symbol in the local scope
		r.SymbolTable[inner] = &localSymbol
	}
	// execute until this top level function returns
	value := r.execUntilReturn()
	// exit the scope
	r.exitScope()
	// return to the original place in the code
	r.ProgramCounter = returnPC
	return value
}

func (e *Expression) IsFunction() bool {
	return e.Operator == BO_FUNCTION_CALL
}

// applyOperator applies the specified operator to the specified values, assuming that the operation has been type checked before
func applyOperator(l *BinaryTypedValue, r *BinaryTypedValue, op BinaryOperator) *BinaryTypedValue {
	switch op {
	case BO_PLUS:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_INT8,
				Value: genericPlus[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_INT16,
				Value: genericPlus[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_INT32,
				Value: genericPlus[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_INT64,
				Value: genericPlus[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_UINT8,
				Value: genericPlus[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_UINT16,
				Value: genericPlus[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_UINT32,
				Value: genericPlus[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_UINT64,
				Value: genericPlus[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BYTE,
				Value: genericPlus[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_FLOAT32,
				Value: genericPlus[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_FLOAT64,
				Value: genericPlus[float64](l.Value, r.Value),
			}
		default:
			panic("invalid type for plus operator")
		}
	case BO_MINUS:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_INT8,
				Value: genericMinus[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_INT16,
				Value: genericMinus[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_INT32,
				Value: genericMinus[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_INT64,
				Value: genericMinus[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_UINT8,
				Value: genericMinus[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_UINT16,
				Value: genericMinus[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_UINT32,
				Value: genericMinus[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_UINT64,
				Value: genericMinus[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BYTE,
				Value: genericMinus[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_FLOAT32,
				Value: genericMinus[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_FLOAT64,
				Value: genericMinus[float64](l.Value, r.Value),
			}
		default:
			panic("invalid type for minus operator")
		}
	case BO_MULTIPLY:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_INT8,
				Value: genericMultiply[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_INT16,
				Value: genericMultiply[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_INT32,
				Value: genericMultiply[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_INT64,
				Value: genericMultiply[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_UINT8,
				Value: genericMultiply[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_UINT16,
				Value: genericMultiply[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_UINT32,
				Value: genericMultiply[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_UINT64,
				Value: genericMultiply[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BYTE,
				Value: genericMultiply[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_FLOAT32,
				Value: genericMultiply[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_FLOAT64,
				Value: genericMultiply[float64](l.Value, r.Value),
			}
		default:
			panic("invalid type for multiply operator")
		}
	case BO_DIVIDE:
		switch l.Type {
		case BT_INT8:
			return &BinaryTypedValue{
				Type:  BT_INT8,
				Value: genericDivide[int8](l.Value, r.Value),
			}
		case BT_INT16:
			return &BinaryTypedValue{
				Type:  BT_INT16,
				Value: genericDivide[int16](l.Value, r.Value),
			}
		case BT_INT32:
			return &BinaryTypedValue{
				Type:  BT_INT32,
				Value: genericDivide[int32](l.Value, r.Value),
			}
		case BT_INT64:
			return &BinaryTypedValue{
				Type:  BT_INT64,
				Value: genericDivide[int64](l.Value, r.Value),
			}
		case BT_UINT8:
			return &BinaryTypedValue{
				Type:  BT_UINT8,
				Value: genericDivide[uint8](l.Value, r.Value),
			}
		case BT_UINT16:
			return &BinaryTypedValue{
				Type:  BT_UINT16,
				Value: genericDivide[uint16](l.Value, r.Value),
			}
		case BT_UINT32:
			return &BinaryTypedValue{
				Type:  BT_UINT32,
				Value: genericDivide[uint32](l.Value, r.Value),
			}
		case BT_UINT64:
			return &BinaryTypedValue{
				Type:  BT_UINT64,
				Value: genericDivide[uint64](l.Value, r.Value),
			}
		case BT_BYTE:
			return &BinaryTypedValue{
				Type:  BT_BYTE,
				Value: genericDivide[byte](l.Value, r.Value),
			}
		case BT_FLOAT32:
			return &BinaryTypedValue{
				Type:  BT_FLOAT32,
				Value: genericDivide[float32](l.Value, r.Value),
			}
		case BT_FLOAT64:
			return &BinaryTypedValue{
				Type:  BT_FLOAT64,
				Value: genericDivide[float64](l.Value, r.Value),
			}
		default:
			panic("invalid type for divide operator")
		}
	default:
		panic("unrecognized operator")
	}
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

func genericPlus[T Numeric](l any, r any) T {
	resL := l.(T)
	resR := r.(T)
	return resL + resR
}

func genericMinus[T Numeric](l any, r any) T {
	resL := l.(T)
	resR := r.(T)
	return resL - resR
}

func genericMultiply[T Numeric](l any, r any) T {
	resL := l.(T)
	resR := r.(T)
	return resL * resR
}

func genericDivide[T Numeric](l any, r any) float64 {
	resL := l.(T)
	resR := r.(T)
	return float64(resL) / float64(resR)
}
