package verifycode_test
import (
    "os"
    "color"
)
func ExampleNewVerifyCode(){
    var verifyCode = NewVerifyCode()
    verifyCode.SetDPI(72)
    verifyCode.SetColor([]color.Color{})
    verifyCode.SetBackground([]color.Color{color.RGBA{255, 255, 255, 255}})
    verifyCode.SetWidthWithHeight(500, 300)
    verifyCode.SetFont([]string{"simsun.ttf"})
    verifyCode.SetFontSize(200)
    verifyCode.SetHinting(false)
    verifyCode.SetKerning(-100, 100)
    osFile, _ := os.Create("Verification code.png")
    defer osFile.Close()
    verifyCode.PNG("abcd", osFile)
}
