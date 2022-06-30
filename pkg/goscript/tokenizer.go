package goscript

import (
	"fmt"
	"regexp"
	"strings"
)

// GSPrimitive is an enum of all primitives in goscript
type GSPrimitive string

const (
	GPS_INT8    GSPrimitive = "int8"
	GSP_INT16   GSPrimitive = "int16"
	GSP_INT32   GSPrimitive = "int32"
	GSP_INT64   GSPrimitive = "int64"
	GSP_UINT8   GSPrimitive = "uint8"
	GSP_UINT16  GSPrimitive = "uint16"
	GSP_UINT32  GSPrimitive = "uint32"
	GSP_UINT64  GSPrimitive = "uint64"
	GSP_STRING  GSPrimitive = "string"
	GSP_CHAR    GSPrimitive = "char"
	GSP_BYTE    GSPrimitive = "byte"
	GSP_FLOAT32 GSPrimitive = "float32"
	GSP_FLOAT64 GSPrimitive = "float64"
	GSP_ANY     GSPrimitive = "any"
)

// iterable list of all primitives
var PRIMITIVES = [...]GSPrimitive{GPS_INT8, GSP_INT16, GSP_INT32, GSP_INT64, GSP_UINT8, GSP_UINT16, GSP_UINT32, GSP_UINT64, GSP_STRING, GSP_CHAR, GSP_BYTE, GSP_FLOAT32, GSP_FLOAT64, GSP_ANY}

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
var SYMBOL_NAME = regexp.MustCompile(`(?m)[a-zA-Z_]{1}[a-zA-Z0-9_]*`)

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
	return SYMBOL_NAME.Match([]byte(token))
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

func (t *Tokenizer) ParseFunction() *FunctionDefinition {
	return nil
}

type TokenizerState byte

// Tokenizer holds all the context required during tokenization of goscript source code
type Tokenizer struct {
	State   TokenizerState
	Context TokenizerContext
}

type TokenizerContext struct {
	funcNameToFunc map[string]*FunctionDefinition
	symbolToIndex  map[string]*int
	indexToType    map[int]*BinaryType
}

// Parse is the main entrypoint for the tokenizer
func (t *Tokenizer) parse(source string) *IntermediateProgram {
	// extract the function definitions from the source code
	functions := t.findFunctions(source)
	// output our function definitions
	fmt.Printf("Funcs:%+v\n", functions)
	// parse the functions
	parsedFunctions := []*FunctionDefinition{}
	for _, function := range functions {
		parsedFunctions = append(parsedFunctions, t.parseFunction(function))
	}
	return nil
}

func (t *Tokenizer) parseFunction(fnc UnparsedFunction) *FunctionDefinition {
	ret := FunctionDefinition{
		Returns: BT_NOTYPE,
	}
	// begin by parsing the functions arguments if any exist
	if len(fnc.Args) > 0 {
		ret.Accepts = t.parseArguments(fnc.Args)
	}
	// parse the return type if the function has one
	if len(fnc.Returns) > 0 {
		ret.Returns = t.parseReturnType(fnc.Returns)
	}
	// finally, parse the body of the function
	ret.Operations = t.parseFunctionBody(fnc.Body)
	return &ret
}

func (t *Tokenizer) parseFunctionBody(body string) []*Operation {
	return nil
}

func (t *Tokenizer) parseTypeToken(token string) IntetmediateType {
	if strings.HasPrefix(clean(token), "[]") {
		return BT_ARRAY
	}
}

func (t *Tokenizer) parseReturnType(returns string) BinaryType {
	switch clean(returns) {
	case "int8":
		return BT_INT8
	case "int16":
		return BT_INT16
	case "int32":
		return BT_INT32
	case "int64":
		return BT_INT64
	case "uint8":
		return BT_UINT8
	case "uint16":
		return BT_UINT16
	case "uint32":
		return BT_UINT32
	case "uint64":
		return BT_UINT64
	case "byte":
		return BT_BYTE
	case "float32":
		return BT_FLOAT32
	case "float64":
		return BT_FLOAT64
	case "string":
		return BT_STRING
	case "bool":
		return BT_BOOLEAN
	case BT_ARRAY:
		// assign the underlying value of value to the underlying value of target
		*target.Value.(*[]*BinaryTypedValue) = *value.Value.(*[]*BinaryTypedValue)
	}
	return BT_NOTYPE
}

func (t *Tokenizer) parseArguments(args string) []*Expression {
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

// clean returns s with all leading and trailing whitespace trimmed
func clean(s string) string {
	return strings.Trim(s, " \n\t\r")
}
