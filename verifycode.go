package verifycode
import (
    "fmt"
    "io"
    "io/ioutil"
    "math/big"
    "crypto/rand"
    "strconv"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "image/jpeg"
    "image/gif"
    "strings"
    "github.com/456vv/verifycode/freetype"
    "github.com/456vv/verifycode/freetype/truetype"
)
//Rander 随机数，返回的随机数是 0-n 的其中一个值。
func Rander(n int) int64 {
    max := big.NewInt(int64(n))
    rnd, err := rand.Int(rand.Reader, max)
    if err != nil {
        return 0
    }
    return rnd.Int64()
}

//RandRange 随机数（范围）
func RandRange(min, max int) int64 {
    var r = Rander(max+1 - min)
    return r + int64(min)
}


//RandomText 随机字符
func RandomText(text string , n int) string {
    var (
        b       []rune
        textL   = len([]rune(text))
        l       int64
        i       int = 1
    )
    stringsReader := strings.NewReader(text)
    for;; i++ {
        l = Rander(textL)
        stringsReader.Seek(l, 0)
        ch, _, err := stringsReader.ReadRune()
        if err != nil {
            continue
        }
        b = append(b, ch)
        if i == n {
            break
        }
    }
    return string(b)
}

//Bounds 边界
type Bounds struct {
    XMin, YMin, XMax, YMax int32                            // X0, Y0, X1, Y1
}

//HMetric 垂直测量
type HMetric struct {
    AdvanceWidth, LeftSideBearing int32                     // 全宽，左跨
}

//VMetric 水平测量
type VMetric struct {
    AdvanceHeight, TopSideBearing int32                     // 全高，上跨
}

//Glyph 字形
type Glyph struct{
    G           *truetype.GlyphBuf                          // 字形对象
    F           *truetype.Font                              // 字体对象
    I           truetype.Index                              // 字体索引
    verifyCode  *VerifyCode                                 // 验证码对象
    text        rune                                        // 字形rune
}

//AdvanceWidth 字形的边界方形宽
func (g *Glyph) AdvanceWidth() int32 {
    return g.G.AdvanceWidth
}

//Width 字形宽
func (g *Glyph) Width() int32 {
    return g.G.B.XMax - g.G.B.XMin
}

//Height 字形高
func (g *Glyph) Height() int32 {
    return g.G.B.YMax - g.G.B.YMin
}

//LeftMargin 字形左跨
func (g *Glyph) LeftMargin() int32 {
    return g.G.B.XMin
}

//TopMargin 字形上跨
func (g *Glyph) TopMargin() int32 {
    return g.G.B.YMin
}

//HMetric 水平测量
func (g *Glyph) HMetric() HMetric {
    return (HMetric)(g.F.HMetric(g.F.FUnitsPerEm(), g.I))
}

//VMetric 垂直测量
func (g *Glyph) VMetric() VMetric {
    return (VMetric)(g.F.VMetric(g.F.FUnitsPerEm(), g.I))
}

//Bounds 字形边界
func (g *Glyph) Bounds() Bounds {
    return (Bounds)(g.G.B)
}

//hinting freetype 的微调字形
func (g *Glyph) hinting() freetype.Hinting {
    if g.verifyCode.Hinting {
        return freetype.FullHinting
    }else{
        return freetype.NoHinting
    }
}

//FontGlyph 字体字形
func (g *Glyph) FontGlyph(size float64, c image.Image) (draw.Image, error) {
    var (
        dx, dy      = int(size), int(size)
        x, y        int
        dpi         = g.verifyCode.DPI
        F           = g.F
    )
    // 新建一个 指定大小的 RGBA位图
    drawImage := image.NewRGBA(image.Rect(0, 0, dx, dy))

    // 画背景
    for y := 0; y < dy; y++ {
        for x := 0; x < dx; x++ {
            // 设置某个点的颜色，依次是 RGBA
            drawImage.Set(x, y, color.RGBA{0, 0, 0, 0})
        }
    }

    //字体转图片
    freetypeContext := freetype.NewContext();
    freetypeContext.SetDPI(dpi)
    freetypeContext.SetClip(drawImage.Bounds())
    freetypeContext.SetFont(F)
    freetypeContext.SetFontSize(size)
    freetypeContext.SetHinting(g.hinting())
    freetypeContext.SetDst(drawImage)
    freetypeContext.SetSrc(c)

    x = 0
    y = int(freetypeContext.PointToFix32(float64(size*0.88))>>8)
    pt := freetype.Pt(x, y)
   _, err := freetypeContext.DrawString(string(g.text), pt)
    return drawImage, err
}



//Font 字体
type Font struct {
    font []*truetype.Font                       // 字体集
}

//OpenFont 字体对象
func NewFont(f []string) (*Font, error) {
    var F []*truetype.Font
    for _, p := range f {
        b, err := ioutil.ReadFile(p)
        if err != nil {
            return nil, err
        }
        font, err := freetype.ParseFont(b)
        if err != nil {
            return nil, err
        }
        F = append(F, font)
    }
    return &Font{
        font: F,
    }, nil
}

//Random 随机字体
func (f *Font) Random() *truetype.Font {
    var (
        fLen = len(f.font)
        n   = Rander(fLen)
    )
    return f.font[n]
}

//Color 颜色集
type Color struct {
    color []color.Color
}

//NewColor 颜色对象
func NewColor(c []string) (*Color, error) {
    var C   []color.Color
    var e   = "NewColor: 颜色%s，被解析错误 >> %s"
    for _, s := range c {
        if len(s) != 9 {
            return nil, fmt.Errorf("NewColor: 十六进制颜色符长度不够 %s", s)
        }
        R, err := strconv.ParseInt(s[1:3], 16, 64)
        if err != nil {
            return nil, fmt.Errorf(e, s, err)
        }
        G, err := strconv.ParseInt(s[3:5], 16, 64)
        if err != nil {
            return nil, fmt.Errorf(e, s, err)
        }
        B, err := strconv.ParseInt(s[5:7], 16, 64)
        if err != nil {
            return nil, fmt.Errorf(e, s, err)
        }
        A, err := strconv.ParseInt(s[7:9], 16, 64)
        if err != nil {
            return nil, fmt.Errorf(e, s, err)
        }
        if (R <= 255 && R >= 0) && (G <= 255 && G >=0) && (B <= 255 && G >= 0) && (A <= 255 && A >= 0) {
            C = append(C, color.RGBA{uint8(R), uint8(G), uint8(B), uint8(A)})
        }else{
            return nil, fmt.Errorf("NewColor: 颜色不正确 %s", s)
        }
    }
    return &Color{
        color: C,
    }, nil
}

//Random 随机颜色
func (c *Color) Random() color.Color {
    var (
        cLen    = len(c.color)
        n       = Rander(cLen)
        colorC  color.Color
    )
    if cLen == 0 {
        //随机生成一个颜色
        colorC = color.RGBA{
            R: uint8(RandRange(128, 255)),
            G: uint8(RandRange(128, 255)),
            B: uint8(RandRange(128, 255)),
            A: uint8(RandRange(128, 255)),
        }
    }else{
        //随机读取列表中的一个颜色
        colorC  = c.color[n]
    }
    return colorC
}

//RandomImage 随机图像
func (c *Color) RandomImage() image.Image {
    return image.NewUniform(c.Random())
}

//VerifyCode 验证码
type VerifyCode struct {
    Width, Height   int                                                               // 宽，高
    DPI             float64                                                           // DPI
    Fonts           *Font                                                            // 字体对象
    Size            float64                                                           // 字体大小
    Colors, Backgrounds    *Color                                                  // 颜色，背景
    Hinting         bool                                                                  // 微调
    KerningMin, KerningMax  int                                                   // 间距
}

//NewVerifyCode 验证码对象
func NewVerifyCode() *VerifyCode {
    return &VerifyCode{
        DPI: 72,
        Size: 12,
           KerningMin: 0,
           KerningMax: 12,
    }
}

//hinting truetype 的微调字形
func (VC *VerifyCode) hinting() truetype.Hinting {
    if VC.Hinting {
        return truetype.FullHinting
    }else{
        return truetype.NoHinting
    }
}

//SetDPI 设置图片的DPI。默认为72DPI
func (VC *VerifyCode) SetDPI(dpi float64) {
    VC.DPI = dpi
}

//SetColor 设置图片中文字颜色，支持多种颜色，颜色随机颜色。默认为随机颜色，颜色RGBA值范围128-255
func (VC *VerifyCode) SetColor(c *Color) {
    VC.Colors = c
}

//SetBackground 设置图片背景颜色，支持多种颜色，颜色随机颜色。默认为透明色
func (VC *VerifyCode) SetBackground(c *Color) {
    VC.Backgrounds = c
}

//SetWidthWithHeight 设置图片的宽和高
func (VC *VerifyCode) SetWidthWithHeight(width, height int) {
    VC.Width, VC.Height = width, height
}

//SetFont 设置图片中验证码字体，支持多种字体，字体是随机选择生成水印
func (VC *VerifyCode) SetFont(font *Font) {
    VC.Fonts = font
}
//SetFontSize 设置字体大小，字体过大会超出。
func (VC *VerifyCode) SetFontSize(size float64) {
    VC.Size = size
}
//SetHinting 设置是否微调字形
func (VC *VerifyCode) SetHinting(h bool) {
    VC.Hinting = h
}

//SetKerning 设置验证码之前的间距
func (VC *VerifyCode) SetKerning(min, max int) {
    VC.KerningMin = min
    VC.KerningMax = max
}

//Glyph 字形
func (VC *VerifyCode) Glyph(s rune) (*Glyph, error) {
    var(
        err     error
        index   truetype.Index
    )
    truetypeFont := VC.Fonts.Random()
    truetypeGlyphBuf := truetype.NewGlyphBuf()
    index = truetypeFont.Index(s)
    err = truetypeGlyphBuf.Load(truetypeFont, truetypeFont.FUnitsPerEm(), index, VC.hinting())
    if err != nil {
        return nil, err
    }
    return &Glyph{
        G: truetypeGlyphBuf,
        F: truetypeFont,
        I: index,
        verifyCode: VC,
        text: s,
    }, nil
}

//Draw 水印
func (VC *VerifyCode) Draw(text string) (draw.Image, error) {
    var(
        glyph   *Glyph
        err     error
        dImage  []draw.Image
        sp      image.Point
        x, y    int
        sizeI   = int(VC.Size)
        textL   = len([]rune(text))
        i, rnd      int
    )
    imageRectangle := image.Rect(0, 0, VC.Width, VC.Height)
    imageRGBA := image.NewRGBA(imageRectangle)
    if len(VC.Backgrounds.color) > 0 {
        for x := 0; x<VC.Width; x++ {
            for y := 0; y<VC.Height; y++ {
       imageRGBA.Set(x, y, VC.Backgrounds.Random())
            }
        }
    }
    for _, v := range text {
        glyph, err = VC.Glyph(v)
        if err != nil {
            return nil, err
        }
        drawImage, err := glyph.FontGlyph(VC.Size, VC.Colors.RandomImage())
        if err != nil {
            return nil, err
        }
        dImage = append(dImage, drawImage)
        x   = (VC.Width/textL)*i
        if x < (VC.Width/2) {
            rnd = int(RandRange(0, VC.KerningMax))
        }else{
            rnd = int(RandRange(VC.KerningMin, 0))
        }
        x   = ^(x+rnd)
        y   = ^int(RandRange(sizeI, VC.Height)) + sizeI
        sp = image.Pt(x, y)
        draw.Draw(imageRGBA, imageRGBA.Bounds(), drawImage, sp, draw.Over)
        i++
    }
    return imageRGBA, nil
}

// PNG
func (VC *VerifyCode) PNG(text string, w io.Writer) error {
    imageImage, err := VC.Draw(text)
    if err != nil {
        return err
    }
    return png.Encode(w, imageImage)
}

//GIF 保存为GIF格式图片
func (VC *VerifyCode) GIF(text string, w io.Writer) error {
    imageImage, err := VC.Draw(text)
    if err != nil {
        return err
    }
    return gif.Encode(w, imageImage, &gif.Options{NumColors: 256})
}

//JPEG 保存为JPEG格式图片
func (VC *VerifyCode) JPEG(text string, w io.Writer) error {
    imageImage, err := VC.Draw(text)
    if err != nil {
        return err
    }
    return jpeg.Encode(w, imageImage, &jpeg.Options{Quality: 100})
}
