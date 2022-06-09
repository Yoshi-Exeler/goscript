package gscompiler

import (
	"testing"
)

func TestResolveExpression(t *testing.T) {
	// this tree represents 5 + 5 * 2
	fiveExpr := &Expression{
		Operator: BO_CONSTANT,
		Value:    uint8(5),
		Type:     BT_UINT8,
	}
	twoExpr := &Expression{
		Operator: BO_CONSTANT,
		Value:    uint8(2),
		Type:     BT_UINT8,
	}
	multExpr := &Expression{
		LeftExpression:  fiveExpr,
		RightExpression: twoExpr,
		Operator:        BO_MULTIPLY,
		Type:            BT_UINT8,
	}
	plusExpr := &Expression{
		LeftExpression:  fiveExpr,
		RightExpression: multExpr,
		Operator:        BO_PLUS,
		Type:            BT_UINT8,
	}
	// now resolve it
	result := plusExpr.Resolve()
	if result.Type != BT_UINT8 {
		t.Fatalf("expression resolution failed, expected type %v but got %v", BT_UINT8, result.Type)
	}
	resV := result.Value.(uint8)
	if resV != 15 {
		t.Fatalf("expression resolution failed, expected result %v but got %v", 15, resV)
	}
}
