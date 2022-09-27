application array

func main() {
    let x: List<str> = []
    x = ["one","two","three"]
    print(x[0])
    print(x[1])
    print(x[2])
    x[0] = "four"
    print(x[0])
    x += "five"
    print(x[4])
}