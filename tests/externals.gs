application externaltest

external db from "https://github.com/Yoshi-Exeler/goscript@branch=master"

import "ext/db"
import "ext/jwt"

// the main function does some stuff
func main() {
    let x: string = "db.Connect"
    let y: string = `db.Connect`
    db.Connect()
}

// this is another comment

struct MyStruct {
    A: string
    B: int64
    C: byte
}

func getMyStruct() => MyStruct {
    return nil
}