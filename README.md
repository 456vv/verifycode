例子：

    package main
    import (
      "fmt"
      "os"
      "image/png"
      "github.com/456vv/verifycode"
    )
    func main(){
        //验证码颜色
        c := []string{"#ff8080FF", "#00ff0000", "#8080c0FD"}
        colors, err := verifycode.NewColor(c)
        if err != nil {
            fmt.Println("NewColor: %v", err)
            os.Exit(-1)
        }
        //验证码背景
        b := []string{"#804040FF"}
        backgrounds, err := verifycode.NewColor(b)
        if err != nil {
            fmt.Println("NewColor: %v", err)
            os.Exit(-1)
        }
        //字体
        f := []string{"0.ttf"}
        fonts, err := verifycode.NewFont(f)
        if err != nil {
            fmt.Println("NewFont: %v", err)
            os.Exit(-1)
        }
        verifyCode := verifycode.NewVerifyCode()
        verifyCode.SetDPI(72)           //也可以不用设置这个
        verifyCode.SetColor(colors)
        verifyCode.SetBackground(backgrounds)
        verifyCode.SetFont(fonts)
        verifyCode.SetWidthWithHeight(500, 200) // 宽500px，高200px
        verifyCode.SetFontSize(200)
        verifyCode.SetHinting(false)    //也可以不用设置这个
        verifyCode.SetKerning(-100, 100)    //随机字距，最小-100，最大100
        file, err := os.Create("tmpTest.png")
        if err != nil {
            fmt.Println("创建文件出错 %v", err)
            os.Exit(-1)
        }
        err = verifyCode.PNG("ABCD", file)
        if err != nil {
            fmt.Println("生成验证码出错 %v", err)
            os.Exit(-1)
        }
    }
