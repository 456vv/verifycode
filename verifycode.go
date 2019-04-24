package verifycode

import(
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"math/big"
	"crypto/rand"
	"strings"
	"image"
	"image/color"
	"image/png"
	"image/jpeg"
	"image/gif"
	"image/draw"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)


//Rand 随机数，返回的随机数是 0-n 的其中一个值。
//	n int 		最大限制的数值
//	int64		返回0-n的随机数
func Rand(n int) int64 {
    if n <= 0 {
        return 0
    }
    max := big.NewInt(int64(n))
    rnd, err := rand.Int(rand.Reader, max)
    if err != nil {
        return 0
    }
    return rnd.Int64()
}

//RandRange 随机数（范围）
//	min, max int	最小值，最大值
//	int64			最小至最大范围的值
func RandRange(min, max int) int64 {
    return int64(min) + Rand(max+1 - min)
}

//RandomText 随机字符
//	text string , n int	字符串，指定长度
//	string				返回n个长度的随机字条串
func RandomText(text string , n int) string {
    var (
        b       []rune
        length  = len([]rune(text))
        r       int64
    )
    stringsReader := strings.NewReader(text)
    for i:=0; i<n; i++ {
        r = Rand(length)
        stringsReader.Seek(r, 0)
        ch, _, err := stringsReader.ReadRune()
        if err != nil {
        	i--
            continue
        }
        b = append(b, ch)
    }
    return string(b)
}



//Color 颜色集
type Color struct {
    color []color.Color
}

//AddHEX 增加十六进制颜色
//	text string		颜色字符串，格式如：#11223344 或 11223344
//	error			错误
func (T *Color) AddHEX(text string) error {
	var l = len(text)
	if text == "" || l < 8 || l > 9 {
		return fmt.Errorf("verifycode: 十六进制颜色符长度不符合 %s, 格式如#11223344 或 11223344", text)
	}
	if l > 8 {
		text = text[1:]
	}
    var e   = "verifycode: 颜色 %s，解析错误是 >> %s"
    R, err := strconv.ParseInt(text[0:2], 16, 64)
    if err != nil {
        return fmt.Errorf(e, text, err)
    }
    G, err := strconv.ParseInt(text[2:4], 16, 64)
    if err != nil {
        return fmt.Errorf(e, text, err)
    }
    B, err := strconv.ParseInt(text[4:6], 16, 64)
    if err != nil {
        return fmt.Errorf(e, text, err)
    }
    A, err := strconv.ParseInt(text[6:8], 16, 64)
    if err != nil {
        return fmt.Errorf(e, text, err)
    }
    if (R <= 255 && R >= 0) && (G <= 255 && G >=0) && (B <= 255 && B >= 0) && (A <= 255 && A >= 0) {
        T.color = append(T.color, color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)})
        return nil
    }
	return fmt.Errorf("verifycode: 颜色不正确 %s", text)
}

//AddRGBA 增加RGBA颜色
//	r, g, b, a uint8	RGBA颜色值0-255
//	error				错误
func (T *Color) AddRGBA(r, g, b, a uint8) error {
    if (r <= 255 && r >= 0) && (g <= 255 && g >=0) && (b <= 255 && b >= 0) && (a <= 255 && a >= 0) {
        T.color = append(T.color, color.RGBA{r, g, b, a})
        return nil
    }
	return fmt.Errorf("verifycode: 参数颜色值仅限于0-255范围")
}

//Random 随机颜色
//	color.Color	颜色，如果没有自定义的颜色，否则内部自动生成一个随机颜色
func (T *Color) Random() color.Color {
    var (
        cLen    = len(T.color)
        n       = Rand(cLen)
    )
    if cLen == 0 {
        //随机生成一个颜色
       return color.RGBA{
            R: uint8(RandRange(128, 255)),
            G: uint8(RandRange(128, 255)),
            B: uint8(RandRange(128, 255)),
            A: uint8(RandRange(128, 255)),
        }
    }
    //随机读取列表中的一个颜色
  	return T.color[n]
}


//Font 字体
type Font struct {
    font []*truetype.Font	// 字体集
}

//AddFile 增加字体文件
//	src string	字体文件路径
//	error		错误
func (T *Font) AddFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
    }
    font, err := truetype.Parse(b)
	if err != nil {
		return err
    }
	T.font = append(T.font, font)
	return nil
}

//Random 随机字体
//	*truetype.Font	字体对象
//	error			错误
func (T *Font) Random() (*truetype.Font, error) {
    var fLen = len(T.font)
    if fLen == 0 {
    	return nil, fmt.Errorf("verifycode: 没有可用的字体")
    }
    n := Rand(fLen)
    f := T.font[n]
    return f, nil
}

//Glyph 字形
type Glyph struct{
	Hinting		font.Hinting								// 微调字形
	Size		float64										// 字形大小
	DPI			float64										// PDI，默认72
}

//FontGlyph 字体字形
//	Font *truetype.Font		字体对象
//	text rune				单个文字
//	c image.Image			颜色
//	draw.Image				字形
//	error					错误
func (T *Glyph) FontGlyph(Font *truetype.Font, text rune, c color.Color) (draw.Image, error) {
	var dx, dy      = int(T.Size), int(T.Size)
    // 新建一个 指定大小的 RGBA位图
    drawImage := image.NewRGBA(image.Rect(0, 0, dx, dy))

	//字体转图片
    freetypeContext := freetype.NewContext();
    freetypeContext.SetClip(drawImage.Bounds())
    if T.DPI > 0 && T.DPI < 72 {
		freetypeContext.SetDPI(T.DPI)
    }
    freetypeContext.SetFont(Font)
	freetypeContext.SetFontSize(T.Size)
    freetypeContext.SetHinting(T.Hinting)
    freetypeContext.SetSrc(image.NewUniform(c))
    freetypeContext.SetDst(drawImage)
	
	perEm := Font.FUnitsPerEm()
	glyphBuf := truetype.GlyphBuf{}
	err := glyphBuf.Load(Font, fixed.Int26_6(perEm), Font.Index(text), T.Hinting)
	if err != nil {
		return nil, err
	}
	
	//x轴是从0开始
	fontWidth := int(glyphBuf.AdvanceWidth)
    x := (int(T.Size)-fontWidth)/2
    
    //y轴是从字高开始
    fontHeight := int(glyphBuf.Bounds.Max.Y+glyphBuf.Bounds.Min.Y)
    y := (int(T.Size)-fontHeight)/2+fontHeight
    
    pt := freetype.Pt(x, y)
	pt, err = freetypeContext.DrawString(string(text), pt)
    return drawImage, err
}

//VerifyCode 验证码
type VerifyCode struct {
    Width, Height   int								// 宽，高
    DPI             float64							// DPI
	Font           	Font							// 字体对象
    Size            float64							// 字体大小
	TextColor, BackgroundColor    Color				// 颜色，背景
    Hinting         font.Hinting					// 微调
    SpaceMin, SpaceMax  int							// 间距
}

func NewVerifyCode() *VerifyCode {
	return &VerifyCode{
		Width:800,
		Height:400,
		Size:200,
	}
}

func backgroundColorBlock(v int) int {
	return int(RandRange(int(float32(v)*0.2), int(float32(v)*0.4)))+1

}

//Draw 水印
func (T *VerifyCode) Draw(text string) (draw.Image, error) {
    //绘制一个框大小，也可以说是一张背景
    imageRectangle := image.Rect(0, 0, T.Width, T.Height)
    imageRGBA := image.NewRGBA(imageRectangle)
    
    //绘制背景颜色
    var(
		bgRandH = backgroundColorBlock(T.Width)
		bgRandV = backgroundColorBlock(T.Height)
    	bgColor color.Color
    	yColors = make(map[int]color.Color)
   	)
   	
    yColors[0]=T.BackgroundColor.Random()
    for x := 0; x<T.Width; x++ {
    	if (x+1)%bgRandH == 0 {
			bgRandH = backgroundColorBlock(T.Width)
			yColors = make(map[int]color.Color)
		    yColors[0]=T.BackgroundColor.Random()
    	}
        for y := 0; y<T.Height; y++ {
        	if bgc, ok := yColors[y]; ok {
        		bgColor = bgc
        	}else if (y+1)%bgRandV == 0  {
				bgRandV = backgroundColorBlock(T.Height)
        		yColors[y]=T.BackgroundColor.Random()
        		bgColor = yColors[y]
        	}else{
        		yColors[y] = bgColor
        	}
        	
			imageRGBA.Set(x, y, bgColor)
        }
    }
    
     var(
		sp      	image.Point
        x, y    	int
        sizeI		= int(T.Size)
        textL		= len([]rune(text))	
        rnd      	int
        i			int
    )
	for _, v := range text {
   	   f, err := T.Font.Random()
   	   if err != nil {
   	   	   //字体错误
   	   	   return nil, err
   	   }
		glyph := Glyph{
			Size:T.Size,
			DPI:T.DPI,
			Hinting:T.Hinting,
		}
        drawImage, err := glyph.FontGlyph(f, v, T.TextColor.Random())
        if err != nil {
            return nil, err
        }
        
        x   = (T.Width/textL)*i
        if x < (T.Width/2) {
            rnd = int(RandRange(0, T.SpaceMax))
        }else{
            rnd = int(RandRange(T.SpaceMin, 0))
        }
        x   = ^(x+rnd)
        y   = ^int(RandRange(sizeI, T.Height)) + sizeI
        sp = image.Pt(x, y)
        
		draw.Draw(imageRGBA, imageRectangle, drawImage, sp, draw.Over)
        i++
    }
	return imageRGBA, nil
}

//PNG 生成PNG图片
//	text string	文本
//	w io.Writer	写入接口
//	error		错误
func (T *VerifyCode) PNG(text string, w io.Writer) error {
    imageImage, err := T.Draw(text)
    if err != nil {
        return err
    }
    return png.Encode(w, imageImage)
}


//GIF 保存为GIF格式图片
//	text string	文本
//	w io.Writer	写入接口
//	error		错误
func (T *VerifyCode) GIF(text string, w io.Writer) error {
    imageImage, err := T.Draw(text)
    if err != nil {
        return err
    }
    return gif.Encode(w, imageImage, &gif.Options{NumColors: 256})
}

//JPEG 保存为JPEG格式图片
//	text string	文本
//	w io.Writer	写入接口
//	error		错误
func (T *VerifyCode) JPEG(text string, w io.Writer) error {
    imageImage, err := T.Draw(text)
    if err != nil {
        return err
    }
    return jpeg.Encode(w, imageImage, &jpeg.Options{Quality: 100})
}










