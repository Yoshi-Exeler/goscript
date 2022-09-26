application function

func main() {
    let b: str = printAndReturn(10)
    print(b)
}   

func printAndReturn(a: u64) => u64 {
    print(a)
    return a
}