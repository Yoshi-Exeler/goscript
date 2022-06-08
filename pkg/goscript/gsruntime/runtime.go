package gsruntime

import (
	"fmt"
	"goscript/pkg/goscript/gscompiler"
)

type Runtime struct {
	Symbols        map[string]*gscompiler.BinarySymbol
	ProgramCounter int
	Running        bool
	ExitCode       int // exit code
}

// reset will reset the state of the runtime
func (r *Runtime) reset() {
	r.ProgramCounter = 0
	r.Symbols = make(map[string]*gscompiler.BinarySymbol)
	r.Running = false
	r.ExitCode = 0
}

// Exec will reset the runtime and then run the specified program until it completes
func (r *Runtime) Exec(program gscompiler.Program) int {
	// completely reset the runtime
	r.reset()
	// enter main execution loop
	r.Running = true
	for {
		r.exec(program)
		if !r.Running {
			// if the program has exited, return its exit code
			return r.ExitCode
		}
	}
}

func (r *Runtime) exit(code int) {
	r.ExitCode = code
	r.Running = false
}

// exec executes the operation currently pointed at by the program counter
func (r *Runtime) exec(program gscompiler.Program) {
	// fetch the next operation from the program
	operation := &program[r.ProgramCounter]
	// call the operation handler for this operation type
	switch operation.Type {
	case gscompiler.ASSIGN_EXPRESSION:
		r.execAssignExpression(operation)
	case gscompiler.DECLARE_SYMBOL:
		r.execDeclareSymbol(operation)
	case gscompiler.CONDITIONAL_BLOCK_ENTER:
		r.execConditionalBlockEnter(operation)
	default:
		fmt.Printf("[GSR] runtime exception, invalid operation type %v", operation.Type)
		r.exit(1)
		return
	}
}

func (r *Runtime) execAssignExpression(operation *gscompiler.BinaryOperation) {
}

// TODO: this might be entirely skippable by replacing all variable declarations with default value assignements
// in IR. This would save us the overhead of determinig the default value of the type and rerouting.
func (r *Runtime) execDeclareSymbol(operation *gscompiler.BinaryOperation) {
	// declaring a symbol simply means assigning the symbol its default value
	// therefore we will determinte the symbols default value and then redirect
	// to execAssignExpression

	// argument 0 is the symbol pointer
	symbol, ok := operation.Args[0].(*gscompiler.BinarySymbol)
	// these safeguards should eventually be removed for increased performance
	if !ok {
		fmt.Printf("[GSR] runtime exception, arg0 missing in declaration at %v", r.ProgramCounter)
		r.exit(1)
	}

	// determine the default value for the type of the symbol
	defaultValue := r.defaultValueOf(symbol.Type)

	// inject the default value argument
	operation.Args[1] = defaultValue

	// redirect to assign expression
	r.execAssignExpression(operation)
}

func (r *Runtime) defaultValueOf(binaryType gscompiler.BinaryType) any {
	switch binaryType {
	case gscompiler.BT_INT8:
		return int8(0)
	case gscompiler.BT_INT16:
		return int16(0)
	case gscompiler.BT_INT32:
		return int32(0)
	case gscompiler.BT_INT64:
		return int64(0)
	case gscompiler.BT_UINT8:
		return uint8(0)
	case gscompiler.BT_UINT16:
		return uint16(0)
	case gscompiler.BT_UINT32:
		return uint32(0)
	case gscompiler.BT_UINT64:
		return uint64(0)
	case gscompiler.BT_STRING:
		return ""
	case gscompiler.BT_CHAR:
		return ""
	case gscompiler.BT_BYTE:
		return byte(0)
	case gscompiler.BT_FLOAT32:
		return float32(0)
	case gscompiler.BT_FLOAT64:
		return float64(0)
	case gscompiler.BT_ANY:
		return nil
	case gscompiler.BT_STRUCT:
		return nil
	default:
		fmt.Printf("[GSR] runtime exception, invalid symbol type %v", binaryType)
		r.exit(1)
		return nil
	}
}

func (r *Runtime) execConditionalBlockEnter(operation *gscompiler.BinaryOperation) {

}
