application function

func main() {
    let b: string = printAndReturn(10)
    print(b)
}   

func printAndReturn(a: uint64) => uint64 {
    print(a)
    return a
}