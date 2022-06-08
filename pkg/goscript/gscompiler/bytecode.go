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
