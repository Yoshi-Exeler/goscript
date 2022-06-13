package goscript

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

/*
	func main() {
		let a: uint8 = 11
	}
*/
func TestAssignConstant(t *testing.T) {
	testProgram := Program{
		Operations: []BinaryOperation{
			NewBindOp(1),
			NewAssignExpressionOp(1, NewConstantExpression(uint8(11), BT_UINT8)),
			NewReturnValueOp(NewVSymbolExpression(1, BT_UINT64))},
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
			NewBindOp(1),
			NewBindOp(2),
			NewAssignExpressionOp(1, NewFunctionExpression(3, BT_UINT8, []*FunctionArgument{})), // assign the return value of the function at pc1 to the symbol 1
			NewAssignExpressionOp(2, NewConstantExpression(uint8(11), BT_UINT8)),                // assign the constant 11 to the local symbol 2
			NewReturnValueOp(NewVSymbolExpression(2, BT_UINT8)),                                 // return the value of the symbol 2
		},
		SymbolTableSize: 10,
	}
	fmt.Println(testProgram.String())
	runtime := NewRuntime()
	runtime.Exec(testProgram)
	if runtime.SymbolTable[1].Value.(uint8) != 11 {
		t.Fatalf("symbol should have been 11 but was %v", runtime.SymbolTable[1].Value.(uint8))
	}
	fmt.Printf("%+v\n", *runtime.SymbolTable[1])
}

/*
	func main() {
		let a: uint8 = getLoopIteratorAfter10()
		return a
	}

	func getLoopIteratorAfter10() => uint8 {
		let b: uint64 = 0
		for let i: uint64 = 0; i < 10; i++ {
			b = i
		}
		return b
	}
*/
/*
	Count Loop Snippet
	0 ENTER_SCOPE                    # enter the loop scope
	1 ASSIGN_EXPRESSION 1 CONST(0)   # i := 0
	2 JUMP_IF_NOT 1 < 10 5           # i < 10, if false jump to 5, exiting the loop
	3 ASSIGN_EXPRESSION 1++          # i++
	... actual loop content ...      # do some stuff
	4 JUMP 2                         # go back to the loop head
	5 EXIT_SCOPE                     # exit the loop scope
*/
func TestLoopAssign(t *testing.T) {
	f, err := os.Create("profile.txt")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	testProgram := Program{
		Operations: []BinaryOperation{
			NewBindOp(1),
			NewAssignExpressionOp(1, NewFunctionExpression(3, BT_UINT64, []*FunctionArgument{})), // let a: uint8 = getLoopIteratorAfter10()
			NewReturnValueOp(NewVSymbolExpression(1, BT_UINT64)),
			NewBindOp(2),
			NewAssignExpressionOp(2, NewConstantExpression(uint64(0), BT_UINT64)), // let b: uint64 = 0
			NewEnterScope(), // enter loop scope
			NewBindOp(3),
			NewAssignExpressionOp(3, NewConstantExpression(uint64(0), BT_UINT64)), // let i: uint64 = 0
			NewJumpIfNotOp(12, &Expression{ // break out of loop if i < 10
				LeftExpression:  NewVSymbolExpression(3, BT_UINT64),
				RightExpression: NewConstantExpression(uint64(1000000000), BT_UINT64),
				Operator:        BO_LESSER,
				Value: &BinaryTypedValue{
					Type:  BT_BOOLEAN,
					Value: false,
				},
				Type: BT_BOOLEAN,
			}),
			NewAssignExpressionOp(2, NewVSymbolExpression(3, BT_UINT64)), // b = i
			NewAssignExpressionOp(3, &Expression{ // i++
				LeftExpression:  NewVSymbolExpression(3, BT_UINT64),
				RightExpression: NewConstantExpression(uint64(1), BT_UINT64),
				Operator:        BO_PLUS,
				Value: &BinaryTypedValue{
					Type:  BT_UINT64,
					Value: 0,
				},
				Type: BT_UINT64,
			}),
			NewJumpOp(8),     // go back to the start of the loop
			NewExitScopeOp(), // exit the loop scope
			NewReturnValueOp(NewVSymbolExpression(2, BT_UINT64)), // return b
		},
		SymbolTableSize: 10,
	}
	fmt.Println(testProgram.String())
	runtime := NewRuntime()
	start := time.Now()
	runtime.Exec(testProgram)
	fmt.Printf("completed in %s\n", time.Since(start))
	fmt.Printf("%+v\n", runtime.SymbolTable)
	fmt.Printf("%+v\n", runtime.SymbolScopeStack)
	if runtime.SymbolTable[1].Value.(uint64) != 999999999 {
		t.Fatalf("symbol should have been 999.999.999 but was %v", runtime.SymbolTable[1].Value.(uint64))
	}
	fmt.Printf("%+v\n", *runtime.SymbolTable[1])
}
