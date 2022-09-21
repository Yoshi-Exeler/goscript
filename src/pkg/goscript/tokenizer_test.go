package goscript

import (
	"fmt"
	"testing"
)

func TestRealizeTypeToken(t *testing.T) {

	typeExpr := realizeTypeToken("List<List<*Tensor<float32>>>", false)
	expectValue(typeExpr.Type, BT_LIST)
	expectValue(typeExpr.SubType.Type, BT_LIST)
	expectValue(typeExpr.SubType.SubType.Type, BT_POINTER)
	expectValue(typeExpr.SubType.SubType.SubType.Type, BT_TENSOR)
	expectValue(typeExpr.SubType.SubType.SubType.SubType.Type, BT_FLOAT32)
}

func TestParseExprAddition(t *testing.T) {

	expr := parseExpression(`5+5`)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	res := *val.Value.(*uint64)
	if res != 10 {
		t.Fatalf("expected 5+5 to be 10 but was %v", res)
	}
}

func TestParseExprAdditionF64(t *testing.T) {

	expr := parseExpression(`5.5+5.5`)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	res := *val.Value.(*float64)
	if res != 11 {
		t.Fatalf("expected 5.5+5.5 to be 10 but was %v", res)
	}
}

func TestParseExprSubtraction(t *testing.T) {

	expr := parseExpression(`10-5`)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	res := *val.Value.(*uint64)
	if res != 5 {
		t.Fatalf("expected 10-5 to be 5 but was %v", res)
	}
}

func TestParseExprSubtractionF64(t *testing.T) {

	expr := parseExpression(`10.5-5.5`)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	res := *val.Value.(*float64)
	if res != 5 {
		t.Fatalf("expected 10.5-5.5 to be 10 but was %v", res)
	}
}

func TestParseExprMultiplication(t *testing.T) {

	expr := parseExpression(`10*10`)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	res := *val.Value.(*uint64)
	if res != 100 {
		t.Fatalf("expected 10*10 to be 100 but was %v", res)
	}
}

func TestParseExprMultiplicationF64(t *testing.T) {

	expr := parseExpression(`10.0*0.5`)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	res := *val.Value.(*float64)
	if res != 5 {
		t.Fatalf("expected 10.0*0.5 to be 5 but was %v", res)
	}
}

func TestParseExprDivision(t *testing.T) {

	expr := parseExpression(`10/2`)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	res := *val.Value.(*float64)
	if res != 5 {
		t.Fatalf("expected 10/2 to be 5 but was %v", res)
	}
}

func TestParseExprDivisionF64(t *testing.T) {

	expr := parseExpression(`10.0/2.0`)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	res := *val.Value.(*float64)
	if res != 5 {
		t.Fatalf("expected 10.0/2.0 to be 5 but was %v", res)
	}
}

func TestParseExprSimpleOrder(t *testing.T) {

	expr := parseExpression(`5+10*10`)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	res := *val.Value.(*uint64)
	if res != 105 {
		t.Fatalf("expected 10*10+5 to be 5 but was %v", res)
	}
}

func TestParseExprSimpleBrackets(t *testing.T) {

	expr := parseExpression(`(5+5)*10-(10-3*2)`)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	res := *val.Value.(*uint64)
	if res != 96 {
		t.Fatalf("expected (5+5)*10-(10-2*2) to be 96 but was %v", res)
	}
}

func TestParseFunctionCall(t *testing.T) {

	expr := parseExpression(`test(5*7+1)`)
	fmt.Printf("%+v\n", expr.Value.Value)
}

func TestParseSymbolExpression(t *testing.T) {

	expr := parseExpression(`myVar`)
	fmt.Printf("%+v\n", expr.Value.Value)
}
