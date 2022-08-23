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

func TestParseExprFunc(t *testing.T) {
	tokenizer := &Tokenizer{}
	expr := tokenizer.parseExpression(`5==5 -      9*2.5 +11.5 +  2/1+ test(5+test2(5)) +1-(5*6)`)
	fmt.Printf("TEST::%v\n", expr)
	rt := NewRuntime()
	val := rt.ResolveExpression(expr)
	fmt.Printf("%+v\n", val.Value.(*uint64))
}
