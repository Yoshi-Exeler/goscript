package gscompiler

import (
	"fmt"
	"testing"
)

func TestAssignSymbol(t *testing.T) {
	testProgram := Program{
		Operations:      []BinaryOperation{*NewAssignExpressionOp(1, NewConstantExpression(uint8(11), BT_INT8))},
		SymbolTableSize: 2,
	}
	runtime := NewRuntime()
	runtime.Exec(testProgram)
	if runtime.SymbolTable[1].Value.(uint8) != 11 {
		t.Fatalf("symbol should have been 11 but was %v", runtime.SymbolTable[1].Value.(uint8))
	}
	fmt.Printf("%+v\n", *runtime.SymbolTable[1])
}
