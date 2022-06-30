package goscript

type IntermediateProgram struct {
	Entrypoint string
	Functions  []*FunctionDefinition
}

type Kind byte

const (
	SINGULAR = 1
	ARRAY    = 2
)

type IntetmediateType struct {
	Kind Kind
	Type BinaryType
}

type IntermediateExpression struct {
	Left     *IntermediateExpression
	Right    *IntermediateExpression
	Operator BinaryOperator
	Value    *BinaryTypedValue // only set when the expression is a constant
	Ref      int
	Args     []*FunctionArgument
}

type PartialExpressionType byte

const (
	OPERATOR PartialExpressionType = 1
	CONSTANT PartialExpressionType = 2
	SYMBOL   PartialExpressionType = 3
)

type PartialExpression struct {
	IsOperator bool // otherwise is assumed to be a value of some kind

}

type FunctionDefinition struct {
	Name       string
	Accepts    []*Expression
	Returns    IntetmediateType
	Operations []*Operation
}

type VSymbol struct {
	Name  string
	Value any
	Type  BinaryType
}

type Operation struct {
	Type OperationType
	Args []any
}
