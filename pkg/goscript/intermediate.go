package goscript

type Application struct {
	Operations []Operation
}

type Operation struct {
}

type VSymbol[T any] struct {
	Type  GSPrimitive
	Name  string
	Value T
}

// STSymbol represents a struct symbol (an instance of a struct type)
type STSymbol struct {
	Name            string
	Type            *StructType
	PrimitiveFields []PrimitiveField
}

// StructType represents the a struct type definition
type StructType struct {
	Name            string
	PrimitiveFields map[string]PrimitiveField // map[field name] field struct
	CompositeFields map[string]CompositeField // map[field name] field struct
}

// PrimitiveFieldType represents a primitive field definition or instance
type PrimitiveField struct {
	Name     string
	Exported bool
	Type     GSPrimitive
	Value    any
}

// CompositeFieldType represents a composite field definition or instance
type CompositeField struct {
	Name     string
	Exported bool
	Type     *StructType
	Value    any
}

type FNSymbol struct{}
