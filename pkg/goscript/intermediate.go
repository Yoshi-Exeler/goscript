package goscript

type IntermediateProgram struct {
	Entrypoint string
	Functions  []*FunctionDefinition
}

type FunctionDefinition struct {
	Name       string
	Accepts    []*VSymbol
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
