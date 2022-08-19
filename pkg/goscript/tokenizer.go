package goscript

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"
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
	FOR             GSKeyword = "for"
	FOREACH         GSKeyword = "foreach"
	LET             GSKeyword = "let"
	FUNC            GSKeyword = "func"
	GSK_RETURN      GSKeyword = "return"
	STRUCT          GSKeyword = "struct"
	CONST           GSKeyword = "const"
	BREAK           GSKeyword = "break"
	CLOSING_BRACKET GSKeyword = "}"
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
var KEYWORDS = [...]GSKeyword{FOR, FOREACH, LET, FUNC, GSK_RETURN, STRUCT, CONST, BREAK, CLOSING_BRACKET}

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
	ret := &IntermediateProgram{}
	// extract the function definitions from the source code
	functions := t.findFunctions(source)
	// output our function definitions
	// parse the functions
	parsedFunctions := []*FunctionDefinition{}
	for _, function := range functions {
		parsedFunctions = append(parsedFunctions, t.parseFunction(function))
	}
	ret.Entrypoint = "main"
	ret.Functions = parsedFunctions
	return ret
}

func (t *Tokenizer) parseFunction(fnc UnparsedFunction) *FunctionDefinition {
	ret := FunctionDefinition{}
	// begin by parsing the functions arguments if any exist
	if len(fnc.Args) > 0 {
		ret.Accepts = t.parseArguments(fnc.Args)
	}
	// parse the return type if the function has one
	if len(fnc.Returns) > 0 {
		ret.Returns = t.parseTypeToken(fnc.Returns, true)
	}
	// finally, parse the body of the function
	ret.Operations = t.parseFunctionBody(fnc.Body)
	ret.Name = fnc.Name
	return &ret
}

func (t *Tokenizer) parseFunctionBody(body string) []*IntermediateOperation {
	ret := []*IntermediateOperation{}
	// parse each line of the function body
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		op := t.parseLine(line)
		ret = append(ret, &op)
	}
	return ret
}

func (t *Tokenizer) parseLine(line string) IntermediateOperation {
	tokens := strings.Split(line, " ")
	// if this line is empty return a no-op which we will delete later in optimization
	if len(tokens) == 0 || len(deleteWhitespace(line)) == 0 {
		return IntermediateOperation{
			Type: IM_NOP,
		}
	}
	// enter keyword parsing mode if the line begins with a keyword
	if isKeyword(tokens[0]) {
		return t.parseKeyWordLine(line, tokens[0])
	}
	// otherwise resolve this line in pure expression mode
	return t.parseExpressionLine(line)
}

func deleteWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func (t *Tokenizer) parseKeyWordLine(line string, keyword string) IntermediateOperation {
	switch keyword {
	case "let":
		return t.parseLetLine(line)
	case "for":
		return t.parseForLine(line)
	case "foreach":
		return t.parseForeachLine(line)
	case "break":
		return t.parseBreakLine(line)
	case "return":
		return t.parseReturnLine(line)
	case "}":
		return t.parseClosingBracketLine(line)
	default:
		panic(fmt.Sprintf("encountered unknown keyword %v", keyword))
	}
}

func (t *Tokenizer) parseExpressionLine(line string) IntermediateOperation {
	return IntermediateOperation{
		Type: IM_EXPRESSION,
		Args: []any{t.parseExpression(line)},
	}
}

func (t *Tokenizer) parseExpression(expr string) *Expression {
	tokens := t.tokenizeExpression(expr)
	for idx, tkn := range tokens {
		fmt.Printf("[%v]:%v\n", idx, tkn)
	}
	return nil
}

func (t *Tokenizer) tokenizeExpression(expr string) []string {
	res := []string{}
	current := ""
	inString := false
	lastChunkIsOp := false
	for i := 0; i < len(expr); i++ {
		// if we read a quote while in string, commit the current string as a chunk
		if string(expr[i]) == `"` && inString {
			current += `"`
			if len(current) > 0 {
				res = append(res, current)
			}
			current = ""
			lastChunkIsOp = false
			inString = false
			continue
		}
		// if we read a quote while not in a string, we enter a string
		if string(expr[i]) == `"` && !inString {
			inString = true
			current += `"`
			lastChunkIsOp = false
			continue
		}
		// if we read a space while not in a string, we commit the current chunk
		if string(expr[i]) == ` ` && !inString {
			if len(current) > 0 {
				res = append(res, current)
				current = ""
			}
			lastChunkIsOp = false
			continue
		}
		// if we read an operator while not in a string, commit the current chunk (max operator length is 2 chars)
		if !inString {
			if charIsOperator(string(expr[i])) {
				if len(current) > 0 {
					res = append(res, current)
				}
				res = append(res, string(expr[i]))
				current = ""
				lastChunkIsOp = true
				continue
			}
			if len(expr)-i+1 > 2 && charIsOperator(string(expr[i])+string(expr[i+1])) {
				if len(current) > 0 {
					res = append(res, current)
				}
				res = append(res, string(expr[i])+string(expr[i+1]))
				current = ""
				i++
				lastChunkIsOp = true
				continue
			}
		}
		// if no condition applies, just append the current char to the current segment
		current += string(expr[i])
	}
	if len(current) > 0 {
		lastChunkIsOp = false
		res = append(res, current)
	}
	if inString {
		log.Fatalf("invalid expression. string was not terminated. expr: %v", expr)
	}
	if lastChunkIsOp {
		log.Fatalf("invalid expression. expression cannot end with an operator. expr: %v", expr)
	}
	return res
}

func charIsOperator(c string) bool {
	ops := []string{"+", "-", "*", "/", "!", "==", "!=", ">", "<", ">=", "<="}
	for _, op := range ops {
		if c == op {
			return true
		}
	}
	return false
}

// matches let myvar: string = "dfg" lines
// G1 is the variable name
// G2 is the variable type
// G3 is the optional value that is assigned
var LET_LINE_REGEX = regexp.MustCompile(`(?m)let ([a-zA-z_]{1}[a-zA-Z0-9-_]*): ([a-zA-Z0-9]*)(?:[ \n]= (.*))?`)

func (t *Tokenizer) parseLetLine(line string) IntermediateOperation {
	ret := IntermediateOperation{
		Type: IM_ASSIGN,
		Args: make([]any, 3),
	}
	// use the let line regex to extract the various components of a let line
	matches := LET_LINE_REGEX.FindStringSubmatch(line)
	if len(matches) != 3 && len(matches) != 4 {
		panic(fmt.Sprintf("could not parse assignement, invalid number of matches (expected 3 || 4, got %v)", len(matches)))
	}
	// assign the symbol name to arg0
	ret.Args[0] = matches[1]
	// get the type of the symbol and assign it to arg1
	parsedType := t.parseTypeToken(matches[2], false)
	ret.Args[1] = parsedType
	// if there is a G3 match save it, otherwise determine the default value for this type
	if len(matches) == 4 {
		ret.Args[2] = t.parseExpression(matches[3])
	} else {
		ret.Args[2] = t.generateDefaultValueForType(parsedType)
	}
	return ret
}

func (t *Tokenizer) generateDefaultValueForType(intermType IntermediateType) IntermediateExpression {
	return IntermediateExpression{}
}

// matches for loop heads 'for i: int32 = 0; i < 10; i++ {'
// G1 is the loop iterator name
// G2 is the loop iterator type
// G3 is the initial value of the iterator
// G4 is the loop condition
// G5 is the loop action (increment or decrement)
var FOR_LINE_REGEX = regexp.MustCompile(`(?mU)for let ([a-zA-z_]{1}[a-zA-Z0-9-_]*): ([a-zA-Z0-9]*) = (.*); (.*);(.*) {`)

func (t *Tokenizer) parseForLine(line string) IntermediateOperation {
	ret := IntermediateOperation{
		Type: IM_FOR,
		Args: make([]any, 5),
	}
	// use the let line regex to extract the various components of a let line
	matches := FOR_LINE_REGEX.FindStringSubmatch(line)
	if len(matches) != 6 {
		panic(fmt.Sprintf("unexpected number of segments in for loop match (expected 6 but got %v)", len(matches)))
	}
	// save the name of the iterator to arg0
	ret.Args[0] = matches[1]
	// parse the iterator type and save it to arg1
	parsedType := t.parseTypeToken(matches[2], false)
	ret.Args[1] = parsedType
	// save the initial value of the iterator to arg3
	parsedExpr := t.parseExpression(matches[3])
	ret.Args[2] = parsedExpr
	// parse the loop condition into an expression
	loopCond := t.parseExpression(matches[4])
	ret.Args[3] = loopCond
	increment := false
	// detect wether we are incrementing or decrementing
	if strings.Contains(matches[5], "++") {
		increment = true
		ret.Args[4] = &increment
	} else if strings.Contains(matches[5], "--") {
		ret.Args[4] = &increment
	} else {
		ret.Args[4] = nil
	}
	// return the op
	return ret
}

// matches foreach loop heads 'foreach element in list {'
// GS1 matches the local loop variable name
// GS2 martches the list being iterated over
var FOREACH_LINE_REGEX = regexp.MustCompile(`(?mU)foreach ([a-zA-z_]{1}[a-zA-Z0-9-_]*) in ([a-zA-z_]{1}[a-zA-Z0-9-_]*) {`)

func (t *Tokenizer) parseForeachLine(line string) IntermediateOperation {
	ret := IntermediateOperation{
		Type: IM_FOREACH,
		Args: make([]any, 2),
	}
	// use the let line regex to extract the various components of a let line
	matches := FOREACH_LINE_REGEX.FindStringSubmatch(line)
	if len(matches) != 3 {
		panic(fmt.Sprintf("unexpected number of segments in foreach loop match (expected 3 but got %v)", len(matches)))
	}
	// save the name of the local iteration symbol to arg0
	ret.Args[0] = matches[1]
	// save the name of the iterable symbol to arg1
	ret.Args[1] = matches[2]
	// yield return
	return ret
}

func (t *Tokenizer) parseBreakLine(line string) IntermediateOperation {
	return IntermediateOperation{
		Type: IM_BREAK,
		Args: []any{},
	}
}

// match return statements 'return test'
// G1 matches the name of the symbol being returned
var RETURN_LINE_REGEX = regexp.MustCompile(`(?m)return( [a-zA-Z0-9+\-*/]*)?$`)

func (t *Tokenizer) parseReturnLine(line string) IntermediateOperation {
	ret := IntermediateOperation{
		Type: IM_RETURN,
		Args: make([]any, 1),
	}
	// use the let line regex to extract the various components of a let line
	fmt.Printf("'%v'\n", line)
	matches := RETURN_LINE_REGEX.FindStringSubmatch(line)
	if len(matches) != 2 && len(matches) != 1 {
		panic(fmt.Sprintf("unexpected number of segments in return match (expected 1 || 2 but got %v)", len(matches)))
	}
	// save the name of the symbol being returned to arg0
	ret.Args[0] = matches[1]
	// yield the op
	return ret
}

func (t *Tokenizer) parseClosingBracketLine(line string) IntermediateOperation {
	return IntermediateOperation{
		Type: IM_CLOSING_BRACKET,
		Args: []any{},
	}
}

func (t *Tokenizer) parseTypeToken(token string, allowNoType bool) IntermediateType {
	var ret IntermediateType
	// first, check composition type modifiers and recusively resolve
	if strings.HasPrefix(token, "[]") {
		trim := strings.TrimPrefix(token, "[]")
		subType := t.parseTypeToken(trim, allowNoType)
		// abort the cascade down the three if we encounter a NOTYPE
		if subType.Type == BT_NOTYPE {
			ret.Type = BT_NOTYPE
			ret.Kind = SINGULAR
			return ret
		}
		ret.SubType = &subType
		ret.Kind = ARRAY
		return ret
	}
	singularType := t.singularTokenToSingular(token)
	if singularType == BT_NOTYPE && !allowNoType {
		panic(fmt.Sprintf("encountered invalid type %v", token))
	}
	ret.Type = singularType
	ret.Kind = SINGULAR
	return ret
}

func (t *Tokenizer) singularTokenToSingular(returns string) BinaryType {
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
	}
	return BT_NOTYPE
}

func (t *Tokenizer) parseArguments(args string) []IntermediateVar {
	ret := []IntermediateVar{}
	// split the args into comma separated list of variables and names
	varsWithNames := strings.Split(args, ",")
	for _, varWithName := range varsWithNames {
		words := strings.Split(varWithName, " ")
		if len(words) != 2 {
			panic(fmt.Sprintf("typed variable token %v has an invalid segment length (expected 2 but got %v)%v", varWithName, len(words), words))
		}
		current := IntermediateVar{
			Name: words[0],
			Type: t.parseTypeToken(words[1], false),
		}
		ret = append(ret, current)
	}
	return ret
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
