package goscript

import "fmt"

type Program struct {
	Operations      []BinaryOperation
	SymbolTableSize int
}

type OperationType byte

const (
	ASSIGN      OperationType = 1 // assign an expression resolution to a symbol
	RETURN      OperationType = 3 // return a value
	CALL        OperationType = 4 // call a function without assigning its return value to anything
	ENTER_SCOPE OperationType = 5 // enters a a new scope
	EXIT_SCOPE  OperationType = 9 // exits the current scope
	JUMP        OperationType = 8 // jumps to the address in arg0
	JUMP_IF     OperationType = 6 // jumps to the address in arg1 if the condition in arg0 is true
	JUMP_IF_NOT OperationType = 7 // jumps to the address in arg1 if the condition in arg0 is false
)

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

/*
	If condition Snippet
	0 JUMP_IF_NOT VS(1) > 10 3         # check the if condition, jump over the block if its false
	1 ENTER_SCOPE                      # enter our scope
	... if block content ...           # do some stuff
	2 EXIT_SCOPE                       # exit our scope
	3 ... some next instruction ...
*/

/*
	if else snippet
	0 ENTER_SCOPE                     # in an if-else, we will definetly enter a block
	1 JUMP_IF_NOT VS(1) > 10 5        # check for the else case, if the main condition is not true, we jump into the else block
	2... if block content ...         # do the stuff for the if branch
	4 JUMP 6                          # jump onto the exit scope at the end
	5... else block content ...       # do the stuff for the else branch
	6 EXIT_SCOPE                      # exit the current scope
	7 ... some next instruction ...
*/

func (p *Program) String() string {
	ret := fmt.Sprintf("BEGIN PROGRAM, %v SYMBOLS\n", len(p.Operations))
	for pc, op := range p.Operations {
		ret += fmt.Sprintf("[%v] %v\n", pc, op.String())
	}
	return ret
}

func (b *BinaryOperation) String() string {
	switch b.Type {
	case ASSIGN:
		return fmt.Sprintf("ASSIGN %v %v", b.Args[0], b.Args[1].(*Expression))
	case RETURN:
		return fmt.Sprintf("RETURN %v", b.Args[0].(*Expression))
	case CALL:
		return fmt.Sprintf("CALL %v", b.Args[0].(*Expression))
	case ENTER_SCOPE:
		return fmt.Sprintf("ENTER_SCOPE")
	case EXIT_SCOPE:
		return fmt.Sprintf("EXIT_SCOPE")
	case JUMP:
		return fmt.Sprintf("JUMP %v", b.Args[0].(int))
	case JUMP_IF:
		return fmt.Sprintf("JUMP_IF %v %v", b.Args[0].(*Expression), b.Args[1].(int))
	case JUMP_IF_NOT:
		return fmt.Sprintf("JUMP_IF_NOT %v %v", b.Args[0].(*Expression), b.Args[1].(int))
	default:
		return "INVALID OP"
	}
}

type BinaryOperation struct {
	Type OperationType // which operation should be performed
	Args []any         // the arguments passed to the operation
}

func NewAssignExpressionOp(symbolRef int, expression *Expression) BinaryOperation {
	return BinaryOperation{
		Type: ASSIGN,
		Args: []any{symbolRef, expression},
	}
}

func NewEnterScope() BinaryOperation {
	return BinaryOperation{
		Type: ENTER_SCOPE,
		Args: []any{},
	}
}

func NewExitScopeOp() BinaryOperation {
	return BinaryOperation{
		Type: EXIT_SCOPE,
		Args: []any{},
	}
}

func NewJumpOp(target int) BinaryOperation {
	return BinaryOperation{
		Type: JUMP,
		Args: []any{target - 1},
	}
}

func NewJumpIfOp(target int, condition *Expression) BinaryOperation {
	return BinaryOperation{
		Type: JUMP_IF,
		Args: []any{condition, target - 1},
	}
}

func NewJumpIfNotOp(target int, condition *Expression) BinaryOperation {
	return BinaryOperation{
		Type: JUMP_IF_NOT,
		Args: []any{condition, target - 1},
	}
}

func NewReturnValueOp(value *Expression) BinaryOperation {
	return BinaryOperation{
		Type: RETURN,
		Args: []any{value},
	}
}

func NewCallFunctionOp(functionExpression *Expression) BinaryOperation {
	return BinaryOperation{
		Type: CALL,
		Args: []any{functionExpression},
	}
}

type BinaryType byte

const (
	BT_INT8    BinaryType = 1
	BT_INT16   BinaryType = 2
	BT_INT32   BinaryType = 3
	BT_INT64   BinaryType = 4
	BT_UINT8   BinaryType = 5
	BT_UINT16  BinaryType = 6
	BT_UINT32  BinaryType = 7
	BT_UINT64  BinaryType = 8
	BT_STRING  BinaryType = 9
	BT_CHAR    BinaryType = 10
	BT_BYTE    BinaryType = 11
	BT_FLOAT32 BinaryType = 12
	BT_FLOAT64 BinaryType = 13
	BT_ANY     BinaryType = 14
	BT_STRUCT  BinaryType = 15
	BT_BOOLEAN BinaryType = 16
)

type BinarySymbol struct {
	Name  string
	Value *BinaryTypedValue
}

type BinaryFunctionCall struct {
	BlockEntry int                 // index of the program to jump to, to begin the function execution
	Args       []*FunctionArgument // symbol map for the arguments of the function call. Maps outside symbols to inside symbols
}

type BinaryOperator byte

const (
	BO_CONSTANT       BinaryOperator = 1
	BO_PLUS           BinaryOperator = 2
	BO_MINUS          BinaryOperator = 3
	BO_MULTIPLY       BinaryOperator = 4
	BO_DIVIDE         BinaryOperator = 5
	BO_FUNCTION_CALL  BinaryOperator = 6 // represents a function that returns a constant
	BO_VSYMBOL        BinaryOperator = 7
	BO_EQUALS         BinaryOperator = 8
	BO_GREATER        BinaryOperator = 9
	BO_LESSER         BinaryOperator = 10
	BO_GREATER_EQUALS BinaryOperator = 11
	BO_LESSER_EQUALS  BinaryOperator = 12
)

type BuiltinFunction byte

const (
	BF_LEN     BuiltinFunction = 1
	BF_PRINT   BuiltinFunction = 2
	BF_PRINTLN BuiltinFunction = 3
	BF_PRINTF  BuiltinFunction = 4
	BF_MIN     BuiltinFunction = 5
	BF_MAX     BuiltinFunction = 6
)

// Expression represents an expression tree.
// If the expression is a constant, the left and right pointers wil be nil, while the
// opertor is BO_CONSTANT. In this case value is non nil and holds the constant value of the expression.
// If the expression is a function, the value will hold the function reference.
type Expression struct {
	LeftExpression  *Expression
	RightExpression *Expression
	Operator        BinaryOperator
	Value           any // only set when the expression is a constant
	Type            BinaryType
}

func (e *Expression) String() string {
	switch e.Operator {
	case BO_CONSTANT:
		return fmt.Sprintf("CONSTANT(%v) ", e.Value)
	case BO_FUNCTION_CALL:
		return fmt.Sprintf("FUNCTION@%v(%v)", e.Value.(*BinaryFunctionCall).BlockEntry, e.Value.(*BinaryFunctionCall).Args)
	case BO_VSYMBOL:
		return fmt.Sprintf("VSYMBOL(%v)", e.Value)
	default:
		return fmt.Sprint("EXPRESSION")
	}
}

func NewVSymbolExpression(symbolRef int, valueType BinaryType) *Expression {
	return &Expression{
		LeftExpression:  nil,
		RightExpression: nil,
		Operator:        BO_VSYMBOL,
		Value:           symbolRef,
		Type:            valueType,
	}
}

func NewFunctionExpression(functionPC int, returnType BinaryType, args []*FunctionArgument) *Expression {
	return &Expression{
		Type:            returnType,
		LeftExpression:  nil,
		RightExpression: nil,
		Operator:        BO_FUNCTION_CALL,
		Value: &BinaryFunctionCall{
			BlockEntry: functionPC,
			Args:       args,
		},
	}
}

func NewConstantExpression(value any, valueType BinaryType) *Expression {
	return &Expression{
		Type:            valueType,
		Value:           value,
		Operator:        BO_CONSTANT,
		LeftExpression:  nil,
		RightExpression: nil,
	}
}

type FunctionArgument struct {
	Expression *Expression
	SymbolRef  int
}

type BinaryTypedValue struct {
	Type  BinaryType
	Value any
}

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func (e *Expression) IsConstant() bool {
	return e.Operator == BO_CONSTANT
}

func (e *Expression) IsFunction() bool {
	return e.Operator == BO_FUNCTION_CALL
}

func (e *Expression) isVSymbol() bool {
	return e.Operator == BO_VSYMBOL
}

func (b *BinaryType) isNumeric() bool {
	switch *b {
	case BT_INT8:
		return true
	case BT_INT16:
		return true
	case BT_INT32:
		return true
	case BT_INT64:
		return true
	case BT_UINT8:
		return true
	case BT_UINT16:
		return true
	case BT_UINT32:
		return true
	case BT_UINT64:
		return true
	case BT_BYTE:
		return true
	case BT_FLOAT32:
		return true
	case BT_FLOAT64:
		return true
	default:
		return false
	}
}

func (b *BinaryType) isIntegerType() bool {
	switch *b {
	case BT_INT8:
		return true
	case BT_INT16:
		return true
	case BT_INT32:
		return true
	case BT_INT64:
		return true
	case BT_UINT8:
		return true
	case BT_UINT16:
		return true
	case BT_UINT32:
		return true
	case BT_UINT64:
		return true
	case BT_BYTE:
		return true
	default:
		return false
	}
}

func (b *BinaryType) isFloatType() bool {
	switch *b {
	case BT_FLOAT32:
		return true
	case BT_FLOAT64:
		return true
	default:
		return false
	}
}
