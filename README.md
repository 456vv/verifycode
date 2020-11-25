# verifycode [![Build Status](https://travis-ci.org/456vv/verifycode.svg?branch=master)](https://travis-ci.org/456vv/verifycode)
golang verifycode, 简单的图形验证码生成。


# **列表：**
```go
func Rand(n int) int64																						// 随机数，返回的随机数是 0-n 的其中一个值。
func RandRange(min, max int) int64																			// 随机数（范围）
func RandomText(text string, n int) string                                                                  // 随机字符
type Color struct {}                                                                                   	// 颜色集
    func (T *Color) AddHEX(text string) error                                                               增加十六进制颜色
    func (T *Color) AddRGBA(r, g, b, a uint8) error                                                         增加RGBA颜色
    func (T *Color) Random() color.Color                                                                    随机颜色
type Font struct {}																						// 字体集
    func (T *Font) AddFile(src string) error                                                                // 增加字体文件
    func (T *Font) Random() (*truetype.Font, error)                                                         // 随机字体
type Glyph struct {   																					// 字形
    Hinting font.Hinting 																					// 微调字形
    Size    float64      																					// 字形大小
    DPI     float64      																					// PDI，默认72
}
    func (T *Glyph) FontGlyph(Font *truetype.Font, text rune, c color.Color) (draw.Image, error)           	// 字体字形
type Style struct{																						// 配色
	Font           	*Font																					// 字体对象
    Size            float64																					// 字体大小
	TextColor, BackgroundColor    *Color																	// 颜色，背景
    Hinting         font.Hinting																			// 微调
    TextSpace  		int																						// 间距
}
type VerifyCode struct {																				// 验证码
    Width, Height              	int          																// 宽，高
    DPI                        	float64      																// DPI
}
	func (T *VerifyCode) Style(s *Style) error																// 配色
    func (T *VerifyCode) Draw(text string) (draw.Image, error)                                              // 水印
    func (T *VerifyCode) GIF(text string, w io.Writer, opt *gif.Options) error								// 保存为GIF格式图片
    func (T *VerifyCode) JPEG(text string, w io.Writer, opt *jpeg.Options) error							// 保存为JPEG格式图片
    func (T *VerifyCode) PNG(text string, w io.Writer) error                                                // 保存为PNG格式图片
```