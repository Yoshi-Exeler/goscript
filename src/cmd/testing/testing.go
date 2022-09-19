package main

import (
	"fmt"

	"github.com/Yoshi-Exeler/goscript/pkg/goscript"
)

func main() {
	testProgram := goscript.Program{
		Operations: []goscript.BinaryOperation{
			goscript.NewBindOp(1, goscript.BT_STRING),
			goscript.NewAssignExpressionOp(1, &goscript.Expression{
				Operator: goscript.BO_BUILTIN_CALL,
				Ref:      int(goscript.BF_INPUTLN),
			}), // assign the constant 11 to the local symbol
			goscript.NewReturnValueOp(&goscript.Expression{
				Operator: goscript.BO_BUILTIN_CALL,
				Ref:      int(goscript.BF_PRINT),
				Args: []*goscript.FunctionArgument{
					&goscript.FunctionArgument{
						Expression: goscript.NewVSymbolExpression(1),
					},
				},
			}),
		},
		SymbolTableSize: 10,
	}
	fmt.Println(testProgram.String())
	runtime := goscript.NewRuntime()
	runtime.Exec(testProgram)
}
