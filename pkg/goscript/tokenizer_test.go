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

func TestFindFunctions(t *testing.T) {
	tokenizer := &Tokenizer{}
	interm := tokenizer.parse(stub1)
	fmt.Printf("%+v\n", interm)
}
