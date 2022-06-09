package gscompiler

type Program []BinaryOperation

type OperationType byte

const (
	ASSIGN_EXPRESSION       OperationType = 1 // assign an expression resolution to a symbol
	DECLARE_SYMBOL          OperationType = 2 // pure declaration of a symbol
	CONDITIONAL_BLOCK_ENTER OperationType = 3 // conditional entry into a block (if-else constructs)
)

type BinaryOperation struct {
	Type OperationType // which operation should be performed
	Args []any         // the arguments passed to the operation
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
	Type  BinaryType
	Value any
}

type BinaryOperator byte

const (
	BO_CONSTANT BinaryOperator = 0
	BO_PLUS     BinaryOperator = 1
	BO_MINUS    BinaryOperator = 2
	BO_MULTIPLY BinaryOperator = 3
	BO_DIVIDE   BinaryOperator = 4
)

// Expression represents an expression tree.
// If the expression is a constant, the left and right pointers wil be nil, while the
// opertor is BO_CONSTANT. In this case value is non nil and holds the constant value of the expression.
type Expression struct {
	LeftExpression  *Expression
	RightExpression *Expression
	Operator        BinaryOperator
	Value           any // only set when the expression is a constant
	Type            BinaryType
}

func (e *Expression) Resolve() *BinaryTypedValue {
	// if the expression is constant, return its value
	if e.IsConstant() {
		return &BinaryTypedValue{
			Type:  e.Type,
			Value: e.Value,
		}
	}
	// otherwise, resolve the left expression
	left := e.LeftExpression.Resolve()
	// then resolve the right expression
	right := e.RightExpression.Resolve()
	// finally apply the operator
	return applyOperator(left, right, e.Operator)
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
