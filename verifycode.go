package verifycode
import (
    "fmt"
    "io"
    "io/ioutil"
    "math/big"
    "crypto/rand"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "image/jpeg"
    "image/gif"
    "strings"
    "github.com/456vv/verifycode/freetype"
    "github.com/456vv/verifycode/freetype/truetype"
    "github.com/456vv/verifycode/freetype/raster"
)

//边界
type Bounds struct {
    XMin, YMin, XMax, YMax int32
}
//水平
type HMetric struct {
    AdvanceWidth, LeftSideBearing int32
}
//垂直
type VMetric struct {
    AdvanceHeight, TopSideBearing int32
}

//字形
type Glyph struct{
    G           *truetype.GlyphBuf
    F           *truetype.Font
    I           truetype.Index
    verifyCode  *VerifyCode
    text        rune
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
//FontGlyph 字形转图片
func (g *Glyph) FontGlyph(size float64, c image.Image) (draw.Image, error) {
    var (
        dx, dy      = int(size), int(size)
        x, y        int
        verifyCode  = g.verifyCode
        dpi         = verifyCode.DPI
        F           = g.F
        pt          raster.Point
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
    pt = freetype.Pt(x, y)
   _, err := freetypeContext.DrawString(string(g.text), pt)
    return drawImage, err
}


//验证码
type VerifyCode struct {
    Width, Height   int
    DPI             float64
    Fonts           []string
    Size            float64
    Colors, Backgrounds []color.Color
    Hinting         bool
    KerningMin, KerningMax      int
}
//验证码对象
func NewVerifyCode() *VerifyCode {
    return &VerifyCode{
        DPI: 72,
        Size: 12,
        KerningMin: 0,
        KerningMax: 20,
    }
}
//SetDPI 设置图片的DPI。默认为72DPI
func (VC *VerifyCode) SetDPI(dpi float64) {
    VC.DPI = dpi
}
//SetColor 设置图片中文字颜色，支持多种颜色，颜色随机颜色。默认为随机颜色，颜色RGBA值范围128-255
func (VC *VerifyCode) SetColor(c []color.Color) {
    VC.Colors = c
}
//SetBackground 设置图片背景颜色，支持多种颜色，颜色随机颜色。默认为透明色
func (VC *VerifyCode) SetBackground(b []color.Color) {
    VC.Backgrounds = b
}
//SetWidthWithHeight 设置图片的宽和高
func (VC *VerifyCode) SetWidthWithHeight(width, height int) {
    VC.Width, VC.Height = width, height
}
//SetFont 设置图片中验证码字体，支持多种字体，字体是随机选择生成水印
func (VC *VerifyCode) SetFont(font []string) {
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
//hinting truetype 的微调字形
func (VC *VerifyCode) hinting() truetype.Hinting {
    if VC.Hinting {
        return truetype.FullHinting
    }else{
        return truetype.NoHinting
    }
}

//Rander 随机数，0-n 范围的数字
func (VC *VerifyCode) Rander(n int) int64 {
    max := big.NewInt(int64(n))
    fontN, err := rand.Int(rand.Reader, max)
    if err != nil {
        return 0
    }
    return fontN.Int64()
}
//RandRange 随机数，min-max 范围的数字
func (VC *VerifyCode) RandRange(min, max int) int64 {
    var r = VC.Rander(max+1 - min)
    return r + int64(min)
}
//randColor 随机选择一个颜色，如果颜色表中没有颜色，否随机生成一个颜色，颜色RGBA值范围128-255
func (VC *VerifyCode) randColor(c []color.Color) image.Image {
    var(
        colorLen = len(c)
        colorC  color.Color
        colorN  int64
    )
    if colorLen <= 0 {
        //随机生成一个颜色
        colorC = color.RGBA{
            R: uint8(VC.RandRange(128, 255)),
            G: uint8(VC.RandRange(128, 255)),
            B: uint8(VC.RandRange(128, 255)),
            A: uint8(VC.RandRange(128, 255)),
        }
    }else{
        //随机读取列表中的一个颜色
        colorN  = VC.Rander(colorLen)
        colorC  = c[colorN]
    }
    return image.NewUniform(colorC)
}
//randOpenFont 随机打开一个字体文件，并返回字体数据。
func (VC *VerifyCode) randOpenFont() ([]byte, error) {
    var (
        fontLen = len(VC.Fonts)
        font    string
    )
    if fontLen <= 0 {
        return nil, fmt.Errorf("没有可用的字体，请设置？x.SetFont([]string{\"ooxx.ttf\", ...})\r\n")
    }
    font = VC.Fonts[VC.Rander(fontLen)]
    return ioutil.ReadFile(font)
}

//Font 随机打开一个字体文件，如果没有字体文件，报错?
func (VC *VerifyCode) Glyph(s rune) (*Glyph, error) {
    var(
        b       []byte
        err     error
        index   truetype.Index
    )
    b, err = VC.randOpenFont()
    if err != nil {
        return nil, err
    }
    truetypeFont, err := freetype.ParseFont(b)
    if err != nil {
        return nil, err
    }
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
//Draw 验证码水印
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
    if len(VC.Backgrounds) > 0 {
        for x := 0; x<VC.Width; x++ {
            for y := 0; y<VC.Height; y++ {
                imageRGBA.Set(x, y, VC.randColor(VC.Backgrounds).At(0, 0))
            }
        }
    }
    for _, v := range text {
        glyph, err = VC.Glyph(v)
        if err != nil {
            return nil, err
        }
        drawImage, err := glyph.FontGlyph(VC.Size, VC.randColor(VC.Colors))
        if err != nil {
            return nil, err
        }
        dImage = append(dImage, drawImage)
        x   = (VC.Width/textL)*i
        if x < (VC.Width/2) {
            rnd = int(VC.RandRange(0, VC.KerningMax))
        }else{
            rnd = int(VC.RandRange(VC.KerningMin, 0))
        }
        x   = ^(x+rnd)
        y   = ^int(VC.RandRange(sizeI, VC.Height)) + sizeI
        sp = image.Pt(x, y)
        draw.Draw(imageRGBA, imageRGBA.Bounds(), drawImage, sp, draw.Over)
        i++
    }
    return imageRGBA, nil
}
//PNG 保存为PNG格式图片
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
//Rnd 随机读取字符
func (VC *VerifyCode) Rnd(text string , num int) string {
    var (
        b       []rune
        textL   = len([]rune(text))
        l       int64
        i       int = 1
    )
    stringsReader := strings.NewReader(text)
    for {
        l = VC.Rander(textL)
        stringsReader.Seek(l, 0)
        ch, _, err := stringsReader.ReadRune()
        b = append(b, ch)
        if err != nil {
            continue
        }
        if i == num {
            break
        }
        i++
    }
    return string(b)
}
