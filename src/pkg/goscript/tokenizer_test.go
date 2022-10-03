package goscript

import (
	"fmt"
	"testing"
)

func TestRealizeTypeToken(t *testing.T) {
	typeExpr := parseTypeWithConstraint("List<List<*Tensor<f32>>>", VALID_TYPE)
	expectValue(typeExpr.Type, BT_LIST)
	expectValue(typeExpr.ValueType.Type, BT_LIST)
	expectValue(typeExpr.ValueType.ValueType.Type, BT_POINTER)
	expectValue(typeExpr.ValueType.ValueType.ValueType.Type, BT_TENSOR)
	expectValue(typeExpr.ValueType.ValueType.ValueType.ValueType.Type, BT_FLOAT32)
}

func TestRealizeInvalidTypeNumeric(t *testing.T) {
	expectPanic(func() {
		_ = parseTypeWithConstraint("Vector<str>", VALID_TYPE)
	})
}

func TestRealizeInvalidTypeUncomposed(t *testing.T) {
	expectPanic(func() {
		_ = parseTypeWithConstraint("Vector<List<f32>>", VALID_TYPE)
	})
}

func TestFindFunctions(t *testing.T) {
	src := `>
func #fn_0_main_main() {
let a: str = "hello world"
#fn_6be61bfd_log_logln(a)
let b: u64 = #fn_b169da74_math_add(5,5)
let c: u64 = #fn_b169da74_math_mult(5,5)
return 0
}
>
func #fn_6be61bfd_log_log(msg: str) {
print(msg)
}
>
func #fn_6be61bfd_log_logln(msg: str) {
#fn_6be61bfd_log_log(msg)
}
>
func #fn_b169da74_math_add(a: u64, b: u64) => u64 {
return a + b
}
>
func #fn_b169da74_math_sub(a: u64, b: u64) => u64 {
return a - b
}
>
func #fn_b169da74_math_mult(a: u64, b: u64) => u64 {
return a * b
}
>
func #fn_b169da74_math_div(a: u64, b: u64) => u64 {
return a / b
}
>`
	funcs, _ := findFunctions(src)
	expectLength(funcs, 7, "should find seven function")
}

func TestGetFunctionArgs(t *testing.T) {
	args := getFunctionArgs("test(a,b)")
	expectLength(args, 2, "when parsing a call with two args, two args should be found")
}

func TestRealizeInvalidTypeWhenAcceptable(t *testing.T) {
	_ = parseTypeWithConstraint("cfsdf>", UNCONSTRAINED)
}

func TestRealizePointerChain(t *testing.T) {
	pointerChain := parseTypeWithConstraint("****u64", VALID_TYPE)
	expectValue(pointerChain.Type, BT_POINTER)
	expectValue(pointerChain.ValueType.Type, BT_POINTER)
	expectValue(pointerChain.ValueType.ValueType.Type, BT_POINTER)
	expectValue(pointerChain.ValueType.ValueType.ValueType.Type, BT_POINTER)
	expectValue(pointerChain.ValueType.ValueType.ValueType.ValueType.Type, BT_UINT64)
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
