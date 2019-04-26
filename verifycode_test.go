package verifycode
	
import(
	"testing"
	"os"
	"image/png"
)


func Test_Rand(t *testing.T){
	for i:=0;i<1000;i++{
		if n := Rand(10); n>10 {
    		t.Fatalf("随机数大于: %v", n)
		}
	}
}

func Test_RandRange(t *testing.T){
    var min, max int = 10, 30
    for i:=0;i<1000;i++ {
    	if n := int(RandRange(min, max)); n < min || n > max {
    		t.Fatalf("随机范围超出: %v", n)
    	}
    }
    
}

func Test_RandomText(t *testing.T) {
    text := "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
    n := 4
    for i:=0;i<1000;i++ {
		if code := RandomText(text , n); len(code) != n {
	    	t.Fatalf("非4位验证码: %v", code)
		}
    }
}


func Test_Color_AddHEX(t *testing.T){
	tests := []struct{
		c	string
		err bool
	}{
		{c:"#11223344", err: false},
		{c:"11223344", err: false},
		{c:"#112233445", err: true},
		{c:"1122334", err: true},
	}
	c := Color{}
	for _, test := range tests {
		if err := c.AddHEX(test.c); (err != nil) != test.err {
			t.Fatalf("格式不正确：%v", err)
		}
	}
}


func Test_Color_AddRGBA(t *testing.T){
	tests := []struct{
		r,g,b,a	uint8
		err bool
	}{
		{r:0, g:0, b:0, a:0, err: false},
		{r:128, g:128, b:128, a:128, err: false},
		{r:255, g:255, b:255, a:255, err: false},
	}
	c := Color{}
	for _, test := range tests {
		if err := c.AddRGBA(test.r, test.g, test.b, test.a); (err != nil) != test.err {
			t.Fatalf("格式不正确：%v", err)
		}
	}
}


func Test_Glyph_FontGlyph(t *testing.T){
	c := Color{}
	f := Font{}
	err := f.AddFile("testdata/0.ttf")
	if err != nil {
		t.Fatalf("字体错误：%v",err)
	}
	
	glyph := Glyph{
		DPI:72,
		Size:500,
	}
	
	drawImage, err := glyph.FontGlyph(f.font[0], 'B', c.Random())
	
	if err != nil {
		t.Fatalf("生成字形出错: %v", err)
	}
	filePNG, err := os.OpenFile("testdata/test.png", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatalf("创建图片文件失败: %v", err)
	}
	defer filePNG.Close()

	err = png.Encode(filePNG, drawImage)
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

}


func Test_NewVerifyCode_PNG(t *testing.T){
	
	filePNG, err := os.OpenFile("testdata/test.png", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatalf("创建图片文件失败: %v", err)
	}
	defer filePNG.Close()
	
	f := Font{}
	err = f.AddFile("testdata/0.ttf")
	if err != nil {
		t.Fatalf("字体错误：%v",err)
	}
	
	verifycode := NewVerifyCode()
	verifycode.Size=200
	verifycode.Font = f
	verifycode.TextSpace=50
	err = verifycode.PNG("验证汉字", filePNG)
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}
}

















