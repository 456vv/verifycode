package verifycode
import (
    "fmt"
    "os"
    "testing"
)
func TestRnd(t *testing.T) {
    var verifyCode = NewVerifyCode()
    s := verifyCode.Rnd("qwertyuiopasdfghjklzxcvbnm", 4)
    fmt.Println(s)
}
func TestPNG(t *testing.T) {
    var verifyCode = NewVerifyCode()
    verifyCode.SetDPI(72)
    verifyCode.SetColor([]color.Color{})
    verifyCode.SetBackground([]color.Color{color.RGBA{255, 255, 255, 255}})
    verifyCode.SetWidthWithHeight(500, 300)
    verifyCode.SetFont([]string{"simkai.ttf"})
    verifyCode.SetFontSize(200)
    verifyCode.SetHinting(false)
    verifyCode.SetKerning(-100, 100)
    osFile, _ := os.Create("Verification code.png")
    defer osFile.Close()
    r := verifyCode.Rnd("1234567890qwertyuiopasdfghjklzxcvbnm", 4)
    verifyCode.PNG(r, osFile)

}
