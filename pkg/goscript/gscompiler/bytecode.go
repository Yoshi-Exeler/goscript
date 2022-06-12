package gscompiler

type Program []BinaryOperation

type OperationType byte

const (
	ASSIGN_EXPRESSION       OperationType = 1 // assign an expression resolution to a symbol
	CONDITIONAL_BLOCK_ENTER OperationType = 2 // conditional entry into a block (if-else constructs)
	RETURN_VALUE            OperationType = 3 // return a value
	CALL_FUNCTION           OperationType = 4 // just call a function
	CLOSE_SCOPE             OperationType = 5 // closing bracket of a scope
)

type BinaryOperation struct {
	Type OperationType // which operation should be performed
	Args []any         // the arguments passed to the operation
}

func NewCloseScopeOp() *BinaryOperation {
	return &BinaryOperation{
		Type: CLOSE_SCOPE,
		Args: []any{},
	}
}

func NewAssignExpressionOp(symbolRef int, expression *Expression) *BinaryOperation {
	return &BinaryOperation{
		Type: ASSIGN_EXPRESSION,
		Args: []any{symbolRef, expression},
	}
}

func NewConditionalBlockEnterOp(condition *Expression, ifBlockPC int, nextOpPC int, elseBlockPC int) *BinaryOperation {
	return &BinaryOperation{
		Type: CONDITIONAL_BLOCK_ENTER,
		Args: []any{condition, ifBlockPC, nextOpPC, elseBlockPC},
	}
}

func NewReturnValueOp(value *Expression) *BinaryOperation {
	return &BinaryOperation{
		Type: RETURN_VALUE,
		Args: []any{value},
	}
}

func NewCallFunctionOp(functionExpression *Expression) *BinaryOperation {
	return &BinaryOperation{
		Type: CALL_FUNCTION,
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
	BO_CONSTANT      BinaryOperator = 1
	BO_PLUS          BinaryOperator = 2
	BO_MINUS         BinaryOperator = 3
	BO_MULTIPLY      BinaryOperator = 4
	BO_DIVIDE        BinaryOperator = 5
	BO_FUNCTION_CALL BinaryOperator = 6 // represents a function that returns a constant
	BO_VSYMBOL       BinaryOperator = 7
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
