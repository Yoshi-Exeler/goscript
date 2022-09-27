
external log from "https://github.com/Yoshi-Exeler/goscript@branch=master"

import "ext/log"
import "std/math"

// the main function does some stuff
func main() {
    let a: str = "hello world"
    log.log(a)
    let b: u64 = math.add(5,5)
    let c: u64 = math.mult(5,5)
    return 0
}