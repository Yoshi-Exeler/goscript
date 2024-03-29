package goscript

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

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

// Parse is the main entrypoint for the tokenizer
func parse(source string) *IntermediateProgram {
	fmt.Println("[GSC][parse] begin parsing")
	start := time.Now()
	ret := &IntermediateProgram{}
	// extract the function definitions from the source code
	fmt.Println("[GSC][findFunctions] begin findFunctions")
	startFindFuncs := time.Now()
	functions, mainFunc := findFunctions(source)
	fmt.Printf("[GSC][STAGE_COMPLETION] findFunctions completed in %v\n", time.Since(startFindFuncs))
	// parse the functions
	fmt.Println("[GSC][parseFunctions] begin parseFunctions")
	startParseFuncs := time.Now()
	parsedFunctions := []*FunctionDefinition{}
	for _, function := range functions {
		if function.Name != "#fn_0_main_main" {
			parsedFunctions = append(parsedFunctions, parseFunction(function))
		}
	}
	ret.Entrypoint = *parseFunction(mainFunc)
	ret.Functions = parsedFunctions
	fmt.Printf("[GSC][STAGE_COMPLETION] parseFunctions completed in %v\n", time.Since(startParseFuncs))
	fmt.Printf("[GSC][STAGE_COMPLETION] parsing completed in %v\n", time.Since(start))
	return ret
}

func parseFunction(fnc UnparsedFunction) *FunctionDefinition {
	fmt.Printf("[GSC][parseFunctions] parsing %v", fnc.Name)
	start := time.Now()
	ret := FunctionDefinition{}
	// begin by parsing the functions arguments if any exist
	if len(fnc.Args) > 0 {
		ret.Accepts = parseArguments(fnc.Args)
	}
	// parse the return type if the function has one
	if len(fnc.Returns) > 0 {
		ret.Returns = parseTypeWithConstraint(fnc.Returns, UNCONSTRAINED)
	}
	// finally, parse the body of the function
	ret.Operations = parseFunctionBody(fnc.Body)
	ret.Name = fnc.Name
	fmt.Printf(" OK %v\n", time.Since(start))
	return &ret
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

func parseFunctionBody(body string) []*IntermediateOperation {
	ret := []*IntermediateOperation{}
	// parse each line of the function body
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		op := parseLine(line)
		ret = append(ret, &op)
	}
	return ret
}

func parseLine(line string) IntermediateOperation {
	line = strings.TrimSpace(line)
	tokens := strings.Split(line, " ")
	// if this line is empty return a no-op which we will delete later in optimization
	if len(tokens) == 0 || len(deleteWhitespace(line)) == 0 {
		return IntermediateOperation{
			Type: IM_NOP,
		}
	}
	// enter keyword parsing mode if the line begins with a keyword
	if isKeyword(tokens[0]) {
		return parseKeyWordLine(line, tokens[0])
	}
	// otherwise resolve this line in pure expression mode
	return parseExpressionLine(line)
}

func deleteWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func parseKeyWordLine(line string, keyword string) IntermediateOperation {
	switch keyword {
	case "let":
		return parseLetLine(line)
	case "for":
		return parseForLine(line)
	case "foreach":
		return parseForeachLine(line)
	case "break":
		return parseBreakLine()
	case "return":
		return parseReturnLine(line)
	case "}":
		return parseClosingBracketLine()
	default:
		panic(fmt.Sprintf("encountered unknown keyword %v", keyword))
	}
}

func parseExpressionLine(line string) IntermediateOperation {
	return IntermediateOperation{
		Type: IM_EXPRESSION,
		Args: []any{parseExpression(line)},
	}
}

func parseExpression(expr string) *Expression {
	// tokenize the expression
	tokens := tokenizeExpression(expr)
	// check if the expression still contains operators
	operatorExists := containsOperator(tokens)
	if !operatorExists {
		// if there is more than one token we panic
		if len(tokens) > 1 {
			panic("invalid expression. unused extra tokens.")
		}
		// if there are no more tokens, we just yield nil
		if len(tokens) == 0 {
			return nil
		}
		// if there is exactly one token, we realize it
		return realizeToken(tokens[0])
	}
	// otherwise build and return the operator tree
	return buildExpressionTree(tokens)
}

type ExpressionTreeNode struct {
	ID       uint64
	Left     uint64
	Right    uint64
	Operator BinaryOperator
	IsChild  bool
}

var PRIORITY_OPS = []string{"*", "/"}

func isPrioritizedOp(op string) bool {
	for _, pop := range PRIORITY_OPS {
		if pop == op {
			return true
		}
	}
	return false
}

func buildExpressionTree(tokens []ExpressionToken) *Expression {
	// map out our tokens by their id
	tokensByID := make(map[uint64]*ExpressionToken)
	for _, token := range tokens {
		token := token
		tokensByID[token.ID] = &token
	}
	// map out our nodes by their id
	nodesByID := make(map[uint64]*ExpressionTreeNode)
	// find the highest token id among the tokens used
	maxTokenID := tokens[len(tokens)-1].ID
	// increment by one to get the first free index
	maxTokenID++
	// save the root node
	var root *ExpressionTreeNode
	// enter the main parser loop
	for {
		// if there is only the root token left we exit
		if len(tokens) == 1 || len(tokens) == 0 {
			break
		}
		// find the next operator in the priority chain
		nextOpIndex := findNextOperator(tokens)
		// create a node instance for it
		node := &ExpressionTreeNode{
			ID:       maxTokenID,
			Left:     tokens[uint64(nextOpIndex-1)].ID,
			Right:    tokens[uint64(nextOpIndex+1)].ID,
			Operator: parseOperator(tokens[nextOpIndex].Value),
		}
		nodesByID[node.ID] = node
		maxTokenID++
		// save this as the root node if required
		root = node
		// replace the relevant tokens with a reference
		newTokens := replaceOperation(tokens, nextOpIndex, node.ID)
		// replace our tokens
		tokens = newTokens
	}
	// generate the expression tree
	return generateExpressionTree(root, nodesByID, tokensByID)
}

func generateExpressionTree(cnode *ExpressionTreeNode, nodes map[uint64]*ExpressionTreeNode, tokens map[uint64]*ExpressionToken) *Expression {
	expression := &Expression{
		Operator: cnode.Operator,
	}
	// if the left node is an op node, recursively resolve it
	if nodes[cnode.Left] != nil {
		expression.LeftExpression = generateExpressionTree(nodes[cnode.Left], nodes, tokens)
	}
	// it the left node is a non-op node, realize it
	if tokens[cnode.Left] != nil {
		ctoken := tokens[cnode.Left]
		expression.LeftExpression = realizeToken(*ctoken)
	}
	// if the right node is an op node, recursively resolve it
	if nodes[cnode.Right] != nil {
		expression.RightExpression = generateExpressionTree(nodes[cnode.Right], nodes, tokens)
	}
	// it the left node is a non-op node, realize it
	if tokens[cnode.Right] != nil {
		ctoken := tokens[cnode.Right]
		expression.RightExpression = realizeToken(*ctoken)
	}
	expression.Value = &BinaryTypedValue{
		Type:  expression.LeftExpression.Value.Type,
		Value: defaultValuePtrOf(expression.LeftExpression.Value.Type),
	}
	return expression
}

type FunctionCallPlaceholder struct {
	Name       string
	SymbolName []string
	Args       []*Expression
}

func realizeFunctionCall(token ExpressionToken) *Expression {
	return &Expression{
		LeftExpression:  nil,
		RightExpression: nil,
		Operator:        BO_FUNCTION_CALL_PLACEHOLDER,
		Value: &BinaryTypedValue{
			Type: BT_NOTYPE,
			Value: &FunctionCallPlaceholder{
				Name: getFunctionName(token.Value),
				Args: parseArgumentExpressions(getFunctionArgs(token.Value)),
			},
		},
	}
}

var FUNCTION_ARG_REGEX = regexp.MustCompile(`(?m)(\(.*\))`)

var ARG_SPLIT_REGEX = regexp.MustCompile(`(?m)[\s,(]+([^,\s()]*)`)

func parseArgumentExpressions(exprs []string) []*Expression {
	res := []*Expression{}
	for _, expr := range exprs {
		res = append(res, parseExpression(expr))
	}
	return res
}

func getFunctionArgs(expr string) []string {
	res := []string{}
	argsMatch := FUNCTION_ARG_REGEX.FindString(expr)
	argMatches := ARG_SPLIT_REGEX.FindAllStringSubmatch(argsMatch, -1)
	for _, match := range argMatches {
		if len(match) != 2 {
			panic("invalid regex match in getFunctionArgs, this should never happen")
		}
		res = append(res, strings.TrimSpace(match[1]))
	}
	return res
}

func getFunctionName(expr string) string {
	split := strings.Split(expr, "(")
	if len(split) == 0 {
		panic(fmt.Sprintf("invalid function call %v", expr))
	}
	return split[0]
}

func stripBrackets(expr string) string {
	return strings.TrimPrefix(strings.TrimSuffix(expr, ")"), "(")
}

func replaceOperation(tokens []ExpressionToken, opIndex int, opID uint64) []ExpressionToken {
	newTokens := []ExpressionToken{}
	for idx, token := range tokens {
		// skip both operands
		if idx == opIndex-1 || idx == opIndex+1 {
			continue
		}
		// replace the actual index with out placeholder
		if idx == opIndex {
			newTokens = append(newTokens, ExpressionToken{
				ID:        opID,
				TokenType: TK_REFERENCE,
			})
			continue
		}
		// copy everything else
		newTokens = append(newTokens, token)
	}
	return newTokens
}

func parseOperator(op string) BinaryOperator {
	switch op {
	case "+":
		return BO_PLUS
	case "-":
		return BO_MINUS
	case "*":
		return BO_MULTIPLY
	case "/":
		return BO_DIVIDE
	case "==":
		return BO_EQUALS
	case ">":
		return BO_GREATER
	case "<":
		return BO_LESSER
	case ">=":
		return BO_GREATER_EQUALS
	case "<=":
		return BO_LESSER_EQUALS
	default:
		panic(fmt.Sprintf("invalid expression. %v is not an operator", op))
	}
}

func findNextOperator(tokens []ExpressionToken) int {
	res := -1
	for i := 0; i < len(tokens); i++ {
		// skip all nodes that arent operators
		if tokens[i].TokenType != TK_OPERATOR {
			continue
		}
		// if we currently dont have an operator use the current one
		if res == -1 {
			res = i
			continue
		}
		// upgrade the operator if it is prioritized but the current one is not
		if !isPrioritizedOp(tokens[res].Value) && isPrioritizedOp(tokens[i].Value) {
			return i
		}
	}
	return res
}

// realizeToken will conver a constant, function call or bracket block into an expression
func realizeToken(token ExpressionToken) *Expression {
	switch token.TokenType {
	case TK_LITERAL:
		return realizeLiteral(token)
	case TK_BRACKET:
		return parseExpression(stripBrackets(token.Value))
	case TK_FUNCTION:
		return realizeFunctionCall(token)
	case TK_STRING:
		str := strings.TrimPrefix(strings.TrimSuffix(token.Value, "\""), "\"")
		return &Expression{
			Operator: BO_CONSTANT,
			Value: &BinaryTypedValue{
				Type:  BT_STRING,
				Value: &str,
			},
		}
	default:
		panic(fmt.Sprintf("invalid expression. cannot realize token %+v", token))
	}
}

func realizeLiteral(token ExpressionToken) *Expression {
	// try to parse the token as an uint64
	u64, err := strconv.ParseUint(token.Value, 10, 64)
	if err == nil {
		return &Expression{
			Value: &BinaryTypedValue{
				Type:  BT_UINT64,
				Value: &u64,
			},
			Operator: BO_CONSTANT,
		}
	}
	i64, err := strconv.ParseInt(token.Value, 10, 64)
	if err == nil {
		return &Expression{
			Value: &BinaryTypedValue{
				Type:  BT_INT64,
				Value: &i64,
			},
			Operator: BO_CONSTANT,
		}
	}
	f64, err := strconv.ParseFloat(token.Value, 64)
	if err == nil {
		return &Expression{
			Value: &BinaryTypedValue{
				Type:  BT_FLOAT64,
				Value: &f64,
			},
			Operator: BO_CONSTANT,
		}
	}
	b, err := strconv.ParseBool(token.Value)
	if err == nil {
		return &Expression{
			Value: &BinaryTypedValue{
				Type:  BT_BOOLEAN,
				Value: &b,
			},
			Operator: BO_CONSTANT,
		}
	}
	// if no other mode matched this must be a symbol
	return &Expression{
		Operator: BO_VSYMBOL_PLACEHOLDER,
		Value: &BinaryTypedValue{
			Type:  BT_NOTYPE,
			Value: token.Value,
		},
	}
}

func containsOperator(tokens []ExpressionToken) bool {
	for _, expToken := range tokens {
		if expToken.TokenType == TK_OPERATOR {
			return true
		}
	}
	return false
}

type ExpressionToken struct {
	ID        uint64
	TokenType TokenType
	Value     string
}

func (e *ExpressionToken) String() string {
	// expString := ""
	// switch e.TokenType {
	// case TK_LITERAL:
	// 	expString = "LITERAL"
	// case TK_STRING:
	// 	expString = "STRING"
	// case TK_BRACKET:
	// 	expString = "BRACKET"
	// case TK_FUNCTION:
	// 	expString = "FUNCTION"
	// case TK_OPERATOR:
	// 	expString = "OPERATOR"
	// case TK_REFERENCE:
	// 	expString = "REFERENCE"
	// }
	if e.TokenType == TK_REFERENCE {
		return "REF"
	}
	return fmt.Sprintf("%v", e.Value)
}

type TokenType byte

const (
	TK_LITERAL   TokenType = 1
	TK_STRING    TokenType = 2
	TK_BRACKET   TokenType = 3
	TK_FUNCTION  TokenType = 4
	TK_OPERATOR  TokenType = 5
	TK_REFERENCE TokenType = 6
)

func tokenizeExpression(expr string) []ExpressionToken {
	res := []ExpressionToken{}
	current := ""
	inString := false
	bracketDepth := 0
	inBrackets := false
	isFunction := false
	lastChunkIsOp := false
	for i := 0; i < len(expr); i++ {
		if string(expr[i]) != `)` && string(expr[i]) != `(` && inBrackets {
			current += string(expr[i])
			lastChunkIsOp = false
			continue
		}
		if string(expr[i]) == `)` && inBrackets {
			// update our bracket depth
			if bracketDepth > 0 {
				bracketDepth--
			} else {
				panic("invalid expression. unbalanced brackets.")
			}
			// if we just hit depth 0, commit the current chunk
			if bracketDepth == 0 {
				current += `)`
				if isFunction {
					if len(current) > 0 {
						res = append(res, ExpressionToken{
							TokenType: TK_FUNCTION,
							Value:     current,
						})
					}
				} else {
					if len(current) > 0 {
						res = append(res, ExpressionToken{
							TokenType: TK_BRACKET,
							Value:     current,
						})
					}
				}
				current = ""
				lastChunkIsOp = false
				inBrackets = false
				isFunction = false
			} else {
				// otherwise just add
				current += `)`
			}
			continue
		}
		if string(expr[i]) == `(` {
			if len(current) > 0 {
				isFunction = true
			}
			bracketDepth++
			if !inBrackets {
				inBrackets = true
			}
			current += `(`
			continue
		}
		// if we read a quote while in string, commit the current string as a chunk
		if string(expr[i]) == `"` && inString {
			current += `"`
			if len(current) > 0 {
				res = append(res, ExpressionToken{
					TokenType: TK_STRING,
					Value:     current,
				})
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
		if string(expr[i]) == ` ` && !inString && !inBrackets {
			if len(current) > 0 {
				res = append(res, ExpressionToken{
					TokenType: TK_LITERAL,
					Value:     current,
				})
				current = ""
			}
			lastChunkIsOp = false
			continue
		}
		// if we read an operator while not in a string, commit the current chunk (max operator length is 2 chars)
		if !inString {
			if charIsOperator(string(expr[i])) {
				if len(current) > 0 {
					res = append(res, ExpressionToken{
						TokenType: TK_LITERAL,
						Value:     current,
					})
				}
				res = append(res, ExpressionToken{
					TokenType: TK_OPERATOR,
					Value:     string(expr[i]),
				})
				current = ""
				lastChunkIsOp = true
				continue
			}
			if len(expr)-i+1 > 2 && charIsOperator(string(expr[i])+string(expr[i+1])) {
				if len(current) > 0 {
					res = append(res, ExpressionToken{
						TokenType: TK_LITERAL,
						Value:     current,
					})
				}
				res = append(res, ExpressionToken{
					TokenType: TK_OPERATOR,
					Value:     string(expr[i]) + string(expr[i+1]),
				})
				current = ""
				i++
				lastChunkIsOp = true
				continue
			}
		}
		// if no condition applies, just append the current char to the current segment
		if !charIsWhitespace(string(expr[i])) && !inString {
			current += string(expr[i])
		}
		if inString {
			current += string(expr[i])
		}
	}
	if len(current) > 0 {
		res = append(res, ExpressionToken{
			TokenType: TK_LITERAL,
			Value:     current,
		})
		lastChunkIsOp = false
	}
	if inString {
		panic(fmt.Sprintf("invalid expression. string was not terminated. expr: %v", expr))
	}
	if lastChunkIsOp {
		panic(fmt.Sprintf("invalid expression. expression cannot end with an operator. expr: %v", expr))
	}
	newTokens := []ExpressionToken{}
	for idx, token := range res {
		token := token
		token.ID = uint64(idx) + 1
		newTokens = append(newTokens, token)
	}
	return newTokens
}

func charIsWhitespace(c string) bool {
	trimmed := strings.TrimSpace(c)
	return len(trimmed) == 0
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
var LET_LINE_REGEX = regexp.MustCompile(`(?m)let ([a-zA-z_]{1}[a-zA-Z0-9-_]*): ([a-zA-Z0-9<>]*)(?:[ \n]= (.*))?`)

func parseLetLine(line string) IntermediateOperation {
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
	parsedType := parseTypeWithConstraint(matches[2], VALID_TYPE)
	ret.Args[1] = parsedType
	// if there is a G3 match save it, otherwise determine the default value for this type
	if len(matches) == 4 {
		ret.Args[2] = parseExpression(matches[3])
	} else {
		ret.Args[2] = generateDefaultValueForType(parsedType)
	}
	return ret
}

func generateDefaultValueForType(intermType IntermediateType) IntermediateExpression {
	return IntermediateExpression{}
}

// matches for loop heads 'for i: int32 = 0; i < 10; i++ {'
// G1 is the loop iterator name
// G2 is the loop iterator type
// G3 is the initial value of the iterator
// G4 is the loop condition
// G5 is the loop action (increment or decrement)
var FOR_LINE_REGEX = regexp.MustCompile(`(?mU)for let ([a-zA-z_]{1}[a-zA-Z0-9-_]*): ([a-zA-Z0-9]*) = (.*); (.*);(.*) {`)

func parseForLine(line string) IntermediateOperation {
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
	parsedType := parseTypeWithConstraint(matches[2], VALID_TYPE)
	ret.Args[1] = parsedType
	// save the initial value of the iterator to arg3
	parsedExpr := parseExpression(matches[3])
	ret.Args[2] = parsedExpr
	// parse the loop condition into an expression
	loopCond := parseExpression(matches[4])
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

func parseForeachLine(line string) IntermediateOperation {
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

func parseBreakLine() IntermediateOperation {
	return IntermediateOperation{
		Type: IM_BREAK,
		Args: []any{},
	}
}

// match return statements 'return test'
// G1 matches the name of the symbol being returned
var RETURN_LINE_REGEX = regexp.MustCompile(`(?m)return\s?(.*)?$`)

func parseReturnLine(line string) IntermediateOperation {
	ret := IntermediateOperation{
		Type: IM_RETURN,
		Args: make([]any, 1),
	}
	// use the let line regex to extract the various components of a let line
	matches := RETURN_LINE_REGEX.FindStringSubmatch(line)
	if len(matches) > 2 {
		panic(fmt.Sprintf("unexpected number of segments in return match (expected 1 but got %v)", len(matches)))
	}
	if len(matches) == 0 {
		return ret
	}
	// save the returned expression to arg0
	if len(deleteWhitespace((matches[1]))) != 0 {
		ret.Args[0] = parseExpression(matches[1])
	}
	// yield the op
	return ret
}

func parseClosingBracketLine() IntermediateOperation {
	return IntermediateOperation{
		Type: IM_CLOSING_BRACKET,
		Args: []any{},
	}
}

func parseSigularType(returns string) BinaryType {
	switch clean(returns) {
	case "i8":
		return BT_INT8
	case "i16":
		return BT_INT16
	case "i32":
		return BT_INT32
	case "i64":
		return BT_INT64
	case "u8":
		return BT_UINT8
	case "u16":
		return BT_UINT16
	case "u32":
		return BT_UINT32
	case "u64":
		return BT_UINT64
	case "byte":
		return BT_BYTE
	case "f32":
		return BT_FLOAT32
	case "f64":
		return BT_FLOAT64
	case "str":
		return BT_STRING
	case "char":
		return BT_CHAR
	case "bool":
		return BT_BOOLEAN
	}
	return BT_NOTYPE
}

func parseArguments(args string) []*IntermediateVar {
	ret := []*IntermediateVar{}
	// split the args into comma separated list of variables and names
	varsWithNames := strings.Split(args, ",")
	for _, varWithName := range varsWithNames {
		rawWords := strings.Split(varWithName, " ")
		words := []string{}
		for _, word := range rawWords {
			if len(deleteWhitespace(word)) != 0 {
				words = append(words, deleteWhitespace(word))
			}
		}
		if len(words) != 2 {
			panic(fmt.Sprintf("typed variable token %v has an invalid segment length (expected 2 but got %v)%v", varWithName, len(words), words))
		}
		current := IntermediateVar{
			Name: words[0],
			Type: parseTypeWithConstraint(words[1], VALID_TYPE),
		}
		ret = append(ret, &current)
	}
	return ret
}

var FUNC_REGEX = regexp.MustCompile(`(?msU)func (#[a-zA-Z_]{1}[a-zA-Z0-9_]*)\((.*)\) (?:=> ([a-zA-Z0-9]*) )?{\n(.*)}\n>`)

func findFunctions(source string) ([]UnparsedFunction, UnparsedFunction) {
	funcs := []UnparsedFunction{}
	mainFunc := UnparsedFunction{}
	matches := FUNC_REGEX.FindAllStringSubmatch(source, -1)
	for _, match := range matches {
		fmt.Printf("[GSC][findFunctions] found %v\n", match[1])
		funcs = append(funcs, UnparsedFunction{
			Name:    match[1],
			Args:    match[2],
			Returns: match[3],
			Body:    match[4],
		})
		if match[1] == "#fn_0_main_main" {
			mainFunc = UnparsedFunction{
				Name:    match[1],
				Args:    match[2],
				Returns: match[3],
				Body:    match[4],
			}
		}
	}
	fmt.Printf("[GSC][findFunctions] %v functions found\n", len(funcs))
	return funcs, mainFunc
}

// clean returns s with all leading and trailing whitespace trimmed
func clean(s string) string {
	return strings.Trim(s, " \n\t\r")
}
