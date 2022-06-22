package goscript

import (
	"fmt"
	"regexp"
	"strings"
)

// GSPrimitive is an enum of all primitives in goscript
type GSPrimitive string

const (
	INT8    GSPrimitive = "int8"
	INT16   GSPrimitive = "int16"
	INT32   GSPrimitive = "int32"
	INT64   GSPrimitive = "int64"
	UINT8   GSPrimitive = "uint8"
	UINT16  GSPrimitive = "uint16"
	UINT32  GSPrimitive = "uint32"
	UINT64  GSPrimitive = "uint64"
	STRING  GSPrimitive = "string"
	CHAR    GSPrimitive = "char"
	BYTE    GSPrimitive = "byte"
	FLOAT32 GSPrimitive = "float32"
	FLOAT64 GSPrimitive = "float64"
	ANY     GSPrimitive = "any"
)

// iterable list of all primitives
var PRIMITIVES = [...]GSPrimitive{INT8, INT16, INT32, INT64, UINT8, UINT16, UINT32, UINT64, STRING, CHAR, BYTE, FLOAT32, FLOAT64, ANY}

type GSKeyword string

const (
	FOR        GSKeyword = "for"
	FOREACH    GSKeyword = "foreach"
	LET        GSKeyword = "let"
	FUNC       GSKeyword = "func"
	GSK_RETURN GSKeyword = "return"
	STRUCT     GSKeyword = "struct"
	CONST      GSKeyword = "const"
	BREAK      GSKeyword = "break"
	// these will be implemented once the compiler generally works
	// EXPORTED GSKeyword = "exported"
	// SWITCH   GSKeyword = "switch"
	// CASE     GSKeyword = "case"
	// DEFAULT  GSKeyword = "default"
	// ASYNC    GSKeyword = "async"
	// AWAIT    GSKeyword = "await"
	// CONTINUE GSKeyword = "continue"
)

// iterable list of all keywords
var KEYWORDS = [...]GSKeyword{FOR, FOREACH, LET, FUNC, GSK_RETURN, STRUCT, CONST, BREAK}

// regex that matches any symbol in goscript
var SYMBOL = regexp.MustCompile(`(?m)[a-zA-Z_]{1}[a-zA-Z0-9_]*`)

// regexes for all primitive literals
var STRING_LITERAL = SIMPLE_STRING_REGEX
var MULTILINE_STRING_LITERAL = MULTILINE_STRING_REGEX
var INTEGER_LITERAL = regexp.MustCompile(`(?m)^\-?[0-9]+$`)
var FLOAT_LITERAL = regexp.MustCompile(`(?m)^\-?[0-9]+\.{1}[0-9]+$`)
var CHAR_LITERAL = regexp.MustCompile(`(?mU)'.{1}'`)
var BOOLEAN_LITERAL = regexp.MustCompile(`(?m)^(?:true)?(?:false)?$`)

// boolean constants
const TRUE = "true"
const FALSE = "false"

type SymbolKind byte

const (
	VSYMBOL  SymbolKind = 1
	STSYMBOL SymbolKind = 2
	FNSYMBOL SymbolKind = 3
)

// getSymbolKind returns the kind of the specified symbol, assuming it is guaranteed to be a symbol
func getSymbolKind(symbol string) SymbolKind {
	if strings.HasPrefix(symbol, "st_") {
		return STSYMBOL
	}
	if strings.HasPrefix(symbol, "fn_") {
		return FNSYMBOL
	}
	return VSYMBOL
}

// isPrimitive checks if a token is a primitive
func isPrimitive(token string) bool {
	for _, primitive := range PRIMITIVES {
		if token == string(primitive) {
			return true
		}
	}
	return false
}

// isSymbol checks if a token is a symbol
func isSymbol(token string) bool {
	return SYMBOL.Match([]byte(token))
}

// isKeyword checks if a token is a keyword
func isKeyword(token string) bool {
	for _, keyword := range KEYWORDS {
		if token == string(keyword) {
			return true
		}
	}
	return false
}

type UnparsedFunction struct {
	Name    string
	Args    string
	Returns string
	Body    string
}

// Tokenizer holds all the context required during tokenization of goscript source code
type Tokenizer struct {
	funcNameToFunc map[string]*FunctionDefinition
	symbolToIndex  map[string]*int
	indexToType    map[int]*BinaryType
}

// Parse is the main entrypoint for the tokenizer
func (t *Tokenizer) parse(source string) *IntermediateProgram {
	// extract the function definitions from the source code
	functions := t.findFunctions(source)
	// output our function definitions
	fmt.Printf("%+v\n", functions)
	//
	return nil
}

func (t *Tokenizer) splitToLines(source string) []string {
	return strings.Split(source, "\n")
}

var FUNC_REGEX = regexp.MustCompile(`(?msU)func ([a-zA-Z_]{1}[a-zA-Z0-9_]*)\((.*)\) (?:=> ([a-zA-Z0-9]*) )?{\n(.*)}`)

func (t *Tokenizer) findFunctions(source string) []UnparsedFunction {
	funcs := []UnparsedFunction{}
	matches := FUNC_REGEX.FindAllStringSubmatch(source, -1)
	for _, match := range matches {
		funcs = append(funcs, UnparsedFunction{
			Name:    match[1],
			Args:    match[2],
			Returns: match[3],
			Body:    match[4],
		})
	}
	return funcs
}

func (t *Tokenizer) parseLineOperation(line string) *Operation {
	return nil
}
