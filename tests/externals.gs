application externaltest

external db from "https://github.com/Yoshi-Exeler/goscript@branch=master"

import "ext/db"
import "ext/jwt"

func main() {
    db.Connect()
}