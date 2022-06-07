application externaltest

external db from "https://github.com/Yoshi-Exeler/goscript@branch=master"

import "ext/db"
import "ext/jwt"

// the main function does some stuff
func main() {
    db.Connect()
}

// this is another comment