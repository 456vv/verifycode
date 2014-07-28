package verifycode
import (
    "testing"
    "os"
    "image/png"
)

func TestRander(t *testing.T){
    n := Rander(10)
    t.Fatalf("随机数: %v", n)
}

func TestRandRange(t *testing.T){
    var min, max int = 10, 30
    n := RandRange(min, max)
    t.Fatalf("随机范围: %v", n)
}

func TestRandomText(t *testing.T) {
    text := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
    n := 4
    code := RandomText(text , n)
    t.Fatalf("4位验证码: %v", code)
}

func TestNewFont(t *testing.T) {
    s := []string{"0.ttf"}
    f, err := NewFont(s)
    if err != nil {
        t.Fatalf("打开字体错误: %v", err)
    }
    if len(f.font) != len(s) {
        t.Fatalf("有些字体无法解析")
    }
}

func TestNewFontRandom(t *testing.T) {
    s := []string{"0.ttf"}
    f, err := NewFont(s)
    if err != nil {
        t.Fatalf("打开字体错误: %v", err)
    }
    font := f.Random()
    _ = font
}

func TestNewColor(t *testing.T) {
    s := []string{"#FFFFFFFF", "#E34D8674"}
    c, err := NewColor(s)
    if err != nil {
        t.Fatalf("解析颜色错误: %v", err)
    }
    if len(c.color) != len(s) {
        t.Fatalf("有些颜色无法解析")
    }
}

func TestNewColorRandom(t *testing.T) {
    s := []string{"#FFFFFFFF", "#E34D8674"}
    c, err := NewColor(s)
    if err != nil {
        t.Fatalf("解析颜色错误: %v", err)
    }
    color := c.Random()
    _ = color
}

func TestNewVerifyCode(t *testing.T){
    //验证码颜色
    c := []string{"#ff8080FF", "#00ff0000", "#8080c0FD"}
    colors, err := NewColor(c)
    if err != nil {
        t.Fatalf("NewColor: %v", err)
    }

    //验证码背景
    b := []string{"#FFFFFFFF"}
    backgrounds, err := NewColor(b)
    if err != nil {
        t.Fatalf("NewColor: %v", err)
    }

    //字体
    f := []string{"0.ttf"}
    fonts, err := NewFont(f)
    if err != nil {
        t.Fatalf("NewFont: %v", err)
    }
    verifycode := NewVerifyCode()
    verifycode.SetDPI(72)           //也可以不用设置这个
    verifycode.SetColor(colors)
    verifycode.SetBackground(backgrounds)
    verifycode.SetFont(fonts)
    verifycode.SetWidthWithHeight(500, 200) // 宽500px，高200px
    verifycode.SetFontSize(200)
    verifycode.SetHinting(false)    //也可以不用设置这个
    verifycode.SetKerning(-100, 100)    //随机字距，最小-100，最大100
    glyph, err := verifycode.Glyph('汉')
    if err != nil {
        t.Fatalf("字形出错 %v", err)
    }
    Image, err := glyph.FontGlyph(verifycode.Size, colors.RandomImage())
    if err != nil {
        t.Fatalf("水印生成出错 %v", err)
    }
    file, err := os.Create("tmpTest.png")
    if err != nil {
        t.Fatalf("创建文件出氏 %v", err)
    }
    png.Encode(file, Image)
}

func TestVercifyCodeDraw(t *testing.T){
    //验证码颜色
    c := []string{"#ff8080FF", "#00ff0000", "#8080c0FD"}
    colors, err := NewColor(c)
    if err != nil {
        t.Fatalf("NewColor: %v", err)
    }

    //验证码背景
    b := []string{"#804040FF"}
    backgrounds, err := NewColor(b)
    if err != nil {
        t.Fatalf("NewColor: %v", err)
    }

    //字体
    f := []string{"0.ttf"}
    fonts, err := NewFont(f)
    if err != nil {
        t.Fatalf("NewFont: %v", err)
    }
    verifycode := NewVerifyCode()
    verifycode.SetDPI(72)           //也可以不用设置这个
    verifycode.SetColor(colors)
    verifycode.SetBackground(backgrounds)
    verifycode.SetFont(fonts)
    verifycode.SetWidthWithHeight(500, 200) // 宽500px，高200px
    verifycode.SetFontSize(200)
    verifycode.SetHinting(false)    //也可以不用设置这个
    verifycode.SetKerning(-100, 100)    //随机字距，最小-100，最大100
    drawImage, err := verifycode.Draw("abcd")
    if err != nil {
        t.Fatalf("生成验证码出错 %v", err)
    }

    file, err := os.Create("tmpTest.png")
    if err != nil {
        t.Fatalf("创建文件出错 %v", err)
    }
    png.Encode(file, drawImage)

}

func TestVercifyCodePNG(t *testing.T){
    //验证码颜色
    c := []string{"#ff8080FF", "#00ff0000", "#8080c0FD"}
    colors, err := NewColor(c)
    if err != nil {
        t.Fatalf("NewColor: %v", err)
    }

    //验证码背景
    b := []string{"#804040FF"}
    backgrounds, err := NewColor(b)
    if err != nil {
        t.Fatalf("NewColor: %v", err)
    }

    //字体
    f := []string{"0.ttf"}
    fonts, err := NewFont(f)
    if err != nil {
        t.Fatalf("NewFont: %v", err)
    }
    verifycode := NewVerifyCode()
    verifycode.SetDPI(72)           //也可以不用设置这个
    verifycode.SetColor(colors)
    verifycode.SetBackground(backgrounds)
    verifycode.SetFont(fonts)
    verifycode.SetWidthWithHeight(500, 200) // 宽500px，高200px
    verifycode.SetFontSize(200)
    verifycode.SetHinting(false)    //也可以不用设置这个
    verifycode.SetKerning(-100, 100)    //随机字距，最小-100，最大100

    file, err := os.Create("tmpTest.png")
    if err != nil {
        t.Fatalf("创建文件出错 %v", err)
    }
    err = verifycode.PNG("ABCD", file)
    if err != nil {
        t.Fatalf("生成验证码出错 %v", err)
    }
}
