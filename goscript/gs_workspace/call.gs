
func main() {
    let a: u64 = 5
    let b: str = "test"
    someFunc(a,b)
}

func someFunc(x: u64, y: str) => u64 {
    println(x)
    println(y)
    return 4
}