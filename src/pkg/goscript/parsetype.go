package goscript

import (
	"fmt"
	"regexp"
)

type TypeConstraint byte

const (
	UNCONSTRAINED TypeConstraint = 1
	VALID_TYPE    TypeConstraint = 2
	UNCOMPOSED    TypeConstraint = 3
	NUMERIC       TypeConstraint = 4
	COMPARABLE    TypeConstraint = 5
)

type Composition struct {
	Type        BinaryType
	Pattern     *regexp.Regexp
	Constraints []TypeConstraint
}

var Compositions = []Composition{
	{
		Type:        BT_LIST,
		Pattern:     regexp.MustCompile(`(?m)^List<(.*)>$`),
		Constraints: []TypeConstraint{VALID_TYPE},
	},
	{
		Type:        BT_VECTOR,
		Pattern:     regexp.MustCompile(`(?m)^Vector<(.*)>$`),
		Constraints: []TypeConstraint{NUMERIC},
	},
	{
		Type:        BT_TENSOR,
		Pattern:     regexp.MustCompile(`(?m)^Tensor<(.*)>$`),
		Constraints: []TypeConstraint{NUMERIC},
	},
	{
		Type:        BT_MAP,
		Pattern:     regexp.MustCompile(`(?m)^Map<(.*)\s*,\s*(.*)>$`),
		Constraints: []TypeConstraint{COMPARABLE, VALID_TYPE},
	},
	{
		Type:        BT_POINTER,
		Pattern:     regexp.MustCompile(`(?m)^\*(.*)$`),
		Constraints: []TypeConstraint{VALID_TYPE},
	},
}

// parseTypeWithConstraint parses infinite depth type compositions into a type tree
func parseTypeWithConstraint(token string, constraint TypeConstraint) IntermediateType {
	// if the token we were called on has length 0, return here. panic if a valid type was required
	if len(token) == 0 {
		if constraint != UNCONSTRAINED {
			panic("type expected but got empty string")
		}
		return IntermediateType{Type: BT_NOTYPE, IsComposed: false}
	}
	// initialize the current token
	currentToken := IntermediateType{}
	// try to parse the top level of the type expression as a composition
	for _, composition := range Compositions {
		match := composition.Pattern.FindAllStringSubmatch(token, -1)
		if len(match) != 0 {
			if constraint != VALID_TYPE && constraint != UNCONSTRAINED {
				panic("expected uncomposed type but found composition")
			}
			currentToken.Type = composition.Type
			currentToken.IsComposed = true
			switch composition.Type {
			case BT_LIST:
				value := parseTypeWithConstraint(composition.Pattern.ReplaceAllString(token, "$1"), composition.Constraints[0])
				currentToken.ValueType = &value
				return currentToken
			case BT_VECTOR:
				value := parseTypeWithConstraint(composition.Pattern.ReplaceAllString(token, "$1"), composition.Constraints[0])
				currentToken.ValueType = &value
				return currentToken
			case BT_TENSOR:
				value := parseTypeWithConstraint(composition.Pattern.ReplaceAllString(token, "$1"), composition.Constraints[0])
				currentToken.ValueType = &value
				return currentToken
			case BT_MAP:
				value := parseTypeWithConstraint(composition.Pattern.ReplaceAllString(token, "$1"), composition.Constraints[0])
				key := parseTypeWithConstraint(composition.Pattern.ReplaceAllString(token, "$2"), composition.Constraints[1])
				currentToken.ValueType = &value
				currentToken.KeyType = &key
				return currentToken
			case BT_POINTER:
				value := parseTypeWithConstraint(composition.Pattern.ReplaceAllString(token, "$1"), composition.Constraints[0])
				currentToken.ValueType = &value
				return currentToken
			default:
				panic(fmt.Sprintf("unexpected type composition %v in parseWithTypeConstraint", composition.Type))
			}
		}
	}
	// if we reached this part of the code we no longer have a composition, so
	// just try to realize the remaining token as a singular
	singularType := parseSigularType(token)
	if singularType == BT_NOTYPE && constraint != UNCONSTRAINED {
		panic(fmt.Sprintf("singular type expected but got '%v' while parsing type token\n", token))
	}
	if constraint == NUMERIC && !singularType.isNumeric() {
		panic(fmt.Sprintf("type was constrained to numeric but parser found non-numeric type %v\n", singularType))
	}
	currentToken.Type = singularType
	currentToken.IsComposed = false
	return currentToken
}
