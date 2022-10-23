

func main() {
    let a: u32 = 5
    let b: str = "test"
    // if the variable identities were calculated cleanly,
    // the inner variables of the add function should not collide with the
    // variables declared above and 10 should just be printed
    println(add(5,5))
}

func add(a: u64, b: u64) => u64 {
    return a + b
}