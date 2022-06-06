application externaltest

external db from "https://github.com/Yoshi-Exeler/goscript@branch=master"

import "ext/db"

func main() {
    db.Connect()
}