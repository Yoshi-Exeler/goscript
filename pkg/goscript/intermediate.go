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

type IntermediateType struct {
	Kind    Kind
	Type    BinaryType
	SubType *IntermediateType
}

type IntermediateVar struct {
	Name string
	Type IntermediateType
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
	Accepts    []IntermediateVar
	Returns    IntermediateType
	Operations []*IntermediateOperation
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

type IntermediateOperationType byte

const (
	IM_NOP             IntermediateOperationType = 0
	IM_ASSIGN          IntermediateOperationType = 1
	IM_FOR             IntermediateOperationType = 2
	IM_CLOSING_BRACKET IntermediateOperationType = 3
	IM_CALL            IntermediateOperationType = 4
	IM_BREAK           IntermediateOperationType = 5
	IM_RETURN          IntermediateOperationType = 6
	IM_FOREACH         IntermediateOperationType = 7
	IM_EXPRESSION      IntermediateOperationType = 8
)

type IntermediateOperation struct {
	Type IntermediateOperationType
	Args []any
}
