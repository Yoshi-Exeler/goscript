package goscript

import (
	"fmt"
	"testing"
)

var stub1 = `
func main() {
let x: string = "db.Connect"
let y: string = "db.Connect"
fn_ad923fa3_db_Connect()
}

func fn_ad923fa3_db_doSomeDBStuff(somevar: uint32) => uint64 {
return 0
}

func fn_ad923fa3_db_Connect() {
return
}

func fn_c40313b4_math_add(a: uint64,b: uint64) => uint64 {
return a+b
}

func fn_e0e93cde_crypto_MakeJWT() {
return
}

func fn_a9d72907_jwt_getSomeJWT() => uint64 {
return 1
}
`

// func TestFindFunctions(t *testing.T) {
// 	tokenizer := &Tokenizer{}
// 	interm := tokenizer.parse(stub1)
// 	fmt.Printf("%+v\n", interm)
// }

// func TestParseExprSimple(t *testing.T) {
// 	tokenizer := &Tokenizer{}
// 	expr := tokenizer.parseExpression(`5==5 -      9*2.5 +11.5   2/1+ "this is a string+-*/" +1`)
// 	rt := NewRuntime()
// 	val := rt.ResolveExpression(expr)
// 	fmt.Printf("%+v\n", val.Value.(*uint64))
// }

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
		t.Fatalf("expected 5+5 to be 10 but was %v", res)
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
