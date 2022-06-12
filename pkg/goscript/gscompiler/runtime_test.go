package gscompiler

import (
	"fmt"
	"testing"
)

/*
	func main() {
		let a: uint8 = 11
	}
*/
func TestAssignConstant(t *testing.T) {
	testProgram := Program{
		Operations:      []BinaryOperation{*NewAssignExpressionOp(1, NewConstantExpression(uint8(11), BT_UINT8))},
		SymbolTableSize: 2,
	}
	runtime := NewRuntime()
	runtime.Exec(testProgram)
	if runtime.SymbolTable[1].Value.(uint8) != 11 {
		t.Fatalf("symbol should have been 11 but was %v", runtime.SymbolTable[1].Value.(uint8))
	}
	fmt.Printf("%+v\n", *runtime.SymbolTable[1])
}

/*
	func main() {
		let a: uint8 = getConst()
	}

	func getConst() => uint8 {
		let b: uint8 = 11
		return b
	}
*/
func TestAssignFunctionReturnValue(t *testing.T) {
	testProgram := Program{
		Operations: []BinaryOperation{
			*NewAssignExpressionOp(1, NewFunctionExpression(1, BT_UINT8, []*FunctionArgument{})), // assign the return value of the function at pc1 to the symbol 1
			*NewAssignExpressionOp(2, NewConstantExpression(uint8(11), BT_UINT8)),                // assign the constant 11 to the local symbol 2
			*NewReturnValueOp(NewVSymbolExpression(2, BT_UINT8)),                                 // return the value of the symbol 2
		},
		SymbolTableSize: 10,
	}
	runtime := NewRuntime()
	runtime.Exec(testProgram)
	if runtime.SymbolTable[1].Value.(uint8) != 11 {
		t.Fatalf("symbol should have been 11 but was %v", runtime.SymbolTable[1].Value.(uint8))
	}
	fmt.Printf("%+v\n", *runtime.SymbolTable[1])
}
