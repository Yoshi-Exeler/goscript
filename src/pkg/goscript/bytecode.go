package goscript

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

type Program struct {
	Operations      []BinaryOperation
	SymbolTableSize int
}

func (p *Program) Encode(out string) {
	f, err := os.OpenFile(out, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	gob.Register(BT_ANY)
	gob.Register(Expression{})
	gob.Register([]*BinaryTypedValue{})
	enc := gob.NewEncoder(f)
	err = enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Program) EncodeBSON(out string) {
	buff, err := bson.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(out, buff, 0600)
	if err != nil {
		log.Fatal(err)
	}
}

type OperationType byte

const (
	ASSIGN       OperationType = 1  // assign an expression resolution to a symbol
	INDEX_ASSIGN OperationType = 2  // assign an expression resolution to a an index of an array symbol
	BIND         OperationType = 3  // binds a symbol to the current scope
	RETURN       OperationType = 4  // return a value
	EXPRESSION   OperationType = 5  // call a function without assigning its return value to anything
	ENTER_SCOPE  OperationType = 6  // enters a a new scope
	EXIT_SCOPE   OperationType = 7  // exits the current scope
	JUMP         OperationType = 8  // jumps to the address in arg0
	JUMP_IF      OperationType = 9  // jumps to the address in arg1 if the condition in arg0 is true
	JUMP_IF_NOT  OperationType = 10 // jumps to the address in arg1 if the condition in arg0 is false
	GROW         OperationType = 11 // grows the array symbol in arg0 by the amount of indices in arg1
	SHRINK       OperationType = 12 // shrinks the array symbol in arg0 by the amount of indices in arg1
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
	Count Loop Snippet
	0 ENTER_SCOPE                    # enter the loop scope
	1 BIND 1 UINT64                  # init and alloc the symbol 1
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
	0 ENTER_SCOPE                     # in an if-else, we will definitely enter a block
	1 JUMP_IF_NOT VS(1) > 10 5        # check for the else case, if the main condition is not true, we jump into the else block
	2... if block content ...         # do the stuff for the if branch
	4 JUMP 6                          # jump onto the exit scope at the end
	5... else block content ...       # do the stuff for the else branch
	6 EXIT_SCOPE                      # exit the current scope
	7 ... some next instruction ...
*/

func (p *Program) String() string {
	ret := fmt.Sprintf("BEGIN PROGRAM, %v INSTRUCTIONS\n", len(p.Operations))
	for pc, op := range p.Operations {
		ret += fmt.Sprintf("[%v] %v\n", pc, op.String())
	}
	return ret
}

func (b *BinaryOperation) String() string {
	switch b.Type {
	case ASSIGN:
		return fmt.Sprintf("ASSIGN SYM(%v) %v", b.Args[0], b.Args[1].(*Expression))
	case INDEX_ASSIGN:
		return fmt.Sprintf("INDEX_ASSIGN SYM(%v) IDX(%v) %v", b.Args[0], b.Args[1], b.Args[2].(*Expression))
	case BIND:
		return fmt.Sprintf("BIND SYM(%v) %v", b.Args[0], b.Args[1].(BinaryType).String())
	case RETURN:
		return fmt.Sprintf("RETURN %v", b.Args[0].(*Expression))
	case EXPRESSION:
		return fmt.Sprintf("EXPRESSION %v", b.Args[0].(*Expression))
	case ENTER_SCOPE:
		return "ENTER_SCOPE"
	case EXIT_SCOPE:
		return "EXIT_SCOPE"
	case JUMP:
		return fmt.Sprintf("JUMP [%v]", b.Args[0].(int))
	case JUMP_IF:
		return fmt.Sprintf("JUMP_IF %v [%v]", b.Args[0].(*Expression), b.Args[1].(int))
	case JUMP_IF_NOT:
		return fmt.Sprintf("JUMP_IF_NOT %v [%v]", b.Args[0].(*Expression), b.Args[1].(int))
	case GROW:
		return fmt.Sprintf("GROW SYM(%v) %v", b.Args[0].(int), b.Args[1].(int))
	case SHRINK:
		return fmt.Sprintf("SHRINK SYM(%v) %v", b.Args[0].(int), b.Args[1].(int))
	default:
		return "INVALID OP"
	}
}

func (bv *BinaryTypedValue) String() string {
	switch bv.Type {
	case BT_INT8:
		// return the derefecrenced value
		return fmt.Sprint(*bv.Value.(*int8))
	case BT_INT16:
		return fmt.Sprint(*bv.Value.(*int16))
	case BT_INT32:
		return fmt.Sprint(*bv.Value.(*int32))
	case BT_INT64:
		return fmt.Sprint(*bv.Value.(*int64))
	case BT_UINT8:
		return fmt.Sprint(*bv.Value.(*uint8))
	case BT_UINT16:
		return fmt.Sprint(*bv.Value.(*uint16))
	case BT_UINT32:
		return fmt.Sprint(*bv.Value.(*uint32))
	case BT_UINT64:
		return fmt.Sprint(*bv.Value.(*uint64))
	case BT_BYTE:
		return fmt.Sprint(*bv.Value.(*byte))
	case BT_FLOAT32:
		return fmt.Sprint(*bv.Value.(*float32))
	case BT_FLOAT64:
		return fmt.Sprint(*bv.Value.(*float64))
	case BT_STRING:
		return fmt.Sprint(*bv.Value.(*string))
	case BT_BOOLEAN:
		return fmt.Sprint(*bv.Value.(*bool))
	case BT_CHAR:
		return fmt.Sprintf("%q", *bv.Value.(*rune))
	case BT_LIST:
		return "[...]"
	case BT_EXPRESSION:
		return fmt.Sprintf(bv.Value.(*Expression).String())
	case BT_NOTYPE:
		return "NOTYPE"
	default:
		panic("unexpected type in stringify typed value")
	}
}

type BinaryOperation struct {
	Type OperationType `json:",omitempty" bson:",omitempty"` // which operation should be performed
	Args []any         `json:",omitempty" bson:",omitempty"` // the arguments passed to the operation
}

func NewAssignExpressionOp(symbolRef int, expression *Expression) BinaryOperation {
	return BinaryOperation{
		Type: ASSIGN,
		Args: []any{symbolRef, expression},
	}
}

func NewIndexAssignOp(symbol int, index *Expression, expression *Expression) BinaryOperation {
	return BinaryOperation{
		Type: INDEX_ASSIGN,
		Args: []any{symbol, index, expression},
	}
}

func NewGrowOperation(symbolRef int, amount int, elemType BinaryType) BinaryOperation {
	return BinaryOperation{
		Type: GROW,
		Args: []any{symbolRef, amount, elemType},
	}
}

func NewShrinkOperation(symbolRef int, amount int) BinaryOperation {
	return BinaryOperation{
		Type: SHRINK,
		Args: []any{symbolRef, amount},
	}
}

func NewEnterScope() BinaryOperation {
	return BinaryOperation{
		Type: ENTER_SCOPE,
		Args: []any{},
	}
}

func NewBindOp(symbolRef int, symbolType BinaryType) BinaryOperation {
	return BinaryOperation{
		Type: BIND,
		Args: []any{symbolRef, symbolType},
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

func NewExpressionOp(expr *Expression) BinaryOperation {
	return BinaryOperation{
		Type: EXPRESSION,
		Args: []any{expr},
	}
}

type BinaryType byte

const (
	BT_INT8       BinaryType = 1
	BT_INT16      BinaryType = 2
	BT_INT32      BinaryType = 3
	BT_INT64      BinaryType = 4
	BT_UINT8      BinaryType = 5
	BT_UINT16     BinaryType = 6
	BT_UINT32     BinaryType = 7
	BT_UINT64     BinaryType = 8
	BT_STRING     BinaryType = 9
	BT_CHAR       BinaryType = 10
	BT_BYTE       BinaryType = 11
	BT_FLOAT32    BinaryType = 12
	BT_FLOAT64    BinaryType = 13
	BT_ANY        BinaryType = 14
	BT_STRUCT     BinaryType = 15
	BT_BOOLEAN    BinaryType = 16
	BT_LIST       BinaryType = 17
	BT_NOTYPE     BinaryType = 18
	BT_EXPRESSION BinaryType = 19
	BT_VECTOR     BinaryType = 20
	BT_TENSOR     BinaryType = 21
	BT_MAP        BinaryType = 22
	BT_POINTER    BinaryType = 23
)

func (b BinaryType) String() string {
	switch b {
	case BT_INT8:
		return "INT8"
	case BT_INT16:
		return "INT16"
	case BT_INT32:
		return "INT32"
	case BT_INT64:
		return "INT64"
	case BT_UINT8:
		return "UINT8"
	case BT_UINT16:
		return "UINT16"
	case BT_UINT32:
		return "UINT32"
	case BT_UINT64:
		return "UINT64"
	case BT_STRING:
		return "STRING"
	case BT_CHAR:
		return "CHAR"
	case BT_BYTE:
		return "BYTE"
	case BT_FLOAT32:
		return "FLOAT32"
	case BT_FLOAT64:
		return "FLOAT64"
	case BT_ANY:
		return "ANY"
	case BT_STRUCT:
		return "STRUCT"
	case BT_BOOLEAN:
		return "BOOLEAN"
	case BT_LIST:
		return "ARRAY"
	default:
		return "invalid type"
	}
}

type BinarySymbol struct {
	Name  string            `json:",omitempty" bson:",omitempty"`
	Value *BinaryTypedValue `json:",omitempty" bson:",omitempty"`
}

type BinaryOperator byte

const (
	BO_CONSTANT                  BinaryOperator = 1
	BO_PLUS                      BinaryOperator = 2
	BO_MINUS                     BinaryOperator = 3
	BO_MULTIPLY                  BinaryOperator = 4
	BO_DIVIDE                    BinaryOperator = 5
	BO_FUNCTION_CALL             BinaryOperator = 6 // represents a function that returns a constant
	BO_VSYMBOL                   BinaryOperator = 7
	BO_EQUALS                    BinaryOperator = 8
	BO_GREATER                   BinaryOperator = 9
	BO_LESSER                    BinaryOperator = 10
	BO_GREATER_EQUALS            BinaryOperator = 11
	BO_LESSER_EQUALS             BinaryOperator = 12
	BO_INDEX_INTO                BinaryOperator = 13 // indexes into an array
	BO_FUNCTION_CALL_PLACEHOLDER BinaryOperator = 14
	BO_VSYMBOL_PLACEHOLDER       BinaryOperator = 15
	BO_BUILTIN_CALL              BinaryOperator = 16
)

func (b BinaryOperator) String() string {
	switch b {
	case BO_PLUS:
		return "+"
	case BO_MINUS:
		return "-"
	case BO_MULTIPLY:
		return "*"
	case BO_DIVIDE:
		return "/"
	case BO_EQUALS:
		return "="
	case BO_GREATER:
		return ">"
	case BO_LESSER:
		return "<"
	case BO_GREATER_EQUALS:
		return ">="
	case BO_LESSER_EQUALS:
		return "<="
	default:
		panic("unknown operator")
	}
}

type BuiltinFunction byte

const (
	BF_LEN       BuiltinFunction = 1
	BF_PRINT     BuiltinFunction = 2
	BF_PRINTLN   BuiltinFunction = 3
	BF_PRINTF    BuiltinFunction = 4
	BF_MIN       BuiltinFunction = 5
	BF_MAX       BuiltinFunction = 6
	BF_INPUT     BuiltinFunction = 7
	BF_INPUTLN   BuiltinFunction = 8
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
)

// Expression represents an expression tree.
// If the expression is a constant, the left and right pointers wil be nil, while the
// opertor is BO_CONSTANT. In this case value is non nil and holds the constant value of the expression.
// If the expression is a function, the value will hold the function reference.
type Expression struct {
	LeftExpression  *Expression         `json:",omitempty" bson:",omitempty"`
	RightExpression *Expression         `json:",omitempty" bson:",omitempty"`
	Operator        BinaryOperator      `json:",omitempty" bson:",omitempty"`
	Value           *BinaryTypedValue   `json:",omitempty" bson:",omitempty"`
	Ref             int                 `json:",omitempty" bson:",omitempty"`
	Args            []*FunctionArgument `json:",omitempty" bson:",omitempty"`
}

func (e *Expression) String() string {
	switch e.Operator {
	case BO_CONSTANT:
		return fmt.Sprintf("CONST(%v)", e.Value.String())
	case BO_FUNCTION_CALL:
		return fmt.Sprintf("FUNC[%v](%v)", e.Ref, e.Args)
	case BO_FUNCTION_CALL_PLACEHOLDER:
		return fmt.Sprintf("FUNCTION_PH[%v](%v)", e.Ref, e.Args)
	case BO_VSYMBOL_PLACEHOLDER:
		return fmt.Sprintf("VSYMBOL_PH[%v](%v)", e.Ref, e.Args)
	case BO_INDEX_INTO:
		return fmt.Sprintf("SYM(%v)[%v]", e.Ref, e.Value.String())
	case BO_VSYMBOL:
		return fmt.Sprintf("SYM(%v)", e.Ref)
	default:
		// if this is not a terminating node, we must recursively travers the expression tree
		expStr := e.LeftExpression.String() + " "
		// append the operator string
		expStr += e.Operator.String()
		// append the right expression string
		expStr += " " + e.RightExpression.String()
		return expStr
	}
}

func NewVSymbolExpression(symbolRef int) *Expression {
	return &Expression{
		LeftExpression:  nil,
		RightExpression: nil,
		Operator:        BO_VSYMBOL,
		Ref:             symbolRef,
		Value: &BinaryTypedValue{
			Value: 0,
		},
	}
}

func NewFunctionExpression(functionPC int, args []*FunctionArgument) *Expression {
	return &Expression{
		LeftExpression:  nil,
		RightExpression: nil,
		Operator:        BO_FUNCTION_CALL,
		Value:           nil,
		Ref:             functionPC,
		Args:            args,
	}
}

func NewConstantExpression(value any, valueType BinaryType) *Expression {
	return &Expression{
		Value: &BinaryTypedValue{
			Value: value,
			Type:  valueType,
		},
		Operator:        BO_CONSTANT,
		LeftExpression:  nil,
		RightExpression: nil,
	}
}

func NewArrayExpression(elements []*BinaryTypedValue) *Expression {
	return &Expression{
		Value: &BinaryTypedValue{
			Value: &elements,
			Type:  BT_LIST,
		},
		Operator:        BO_CONSTANT,
		LeftExpression:  nil,
		RightExpression: nil,
	}
}

func NewIndexIntoExpression(symbol int, index *Expression) *Expression {
	return &Expression{
		LeftExpression:  nil,
		RightExpression: nil,
		Operator:        BO_INDEX_INTO,
		Value: &BinaryTypedValue{
			Type:  BT_EXPRESSION,
			Value: index,
		},
		Ref: symbol,
	}
}

type FunctionArgument struct {
	Expression *Expression `json:",omitempty" bson:",omitempty"`
	SymbolRef  int         `json:",omitempty" bson:",omitempty"`
}

type BinaryTypedValue struct {
	Type  BinaryType `json:",omitempty" bson:",omitempty"`
	Value any        `json:",omitempty" bson:",omitempty"`
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
