package verifycode
import (
    "fmt"
)

func ExampleRand(){
    n := Rand(10)
    fmt.Println(n)
}

func ExampleRandRange(){
    var min, max int = 10, 30
    n := RandRange(min, max)
    fmt.Println(n)
}

func ExampleRandomText(){
    text := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
    n := 4
    code := RandomText(text , n)
    fmt.Println(code)
}