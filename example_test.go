package verifycode_test
import (
    "fmt"
)

func ExampleRander(){
    n := Rander(10)
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

func ExampleNewFont() {
    s := []string{"0.ttf"}
    f, err := NewFont(s)
    fmt.Println(f, err)
}

func ExampleNewColor() {
    s := []string{"#FFFFFFFF", "#E34D86FF"}
    c, err := NewColor(s)
    fmt.Println(c, err)
}
