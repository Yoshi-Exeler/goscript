package gscompiler

import "regexp"

// list of primitive types available in goscript
var PRIMITIVES = []string{"int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "string", "char", "byte", "float32", "float64", "any"}

// regex that matches any symbol in goscript
var SYMBOL = regexp.MustCompile(`(?m)[a-zA-Z_]{1}[a-zA-Z0-9_]*`)

// boolean constants
const TRUE = "true"
const FALSE = "false"

func isPrimitive(expr string) bool {
	for _, primitive := range PRIMITIVES {
		if expr == primitive {
			return true
		}
	}
	return false
}
