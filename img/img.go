package img

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
)

type Img struct {
	ImgTpl  image.Image
	ext     string
	ImgOut  *image.NRGBA
	ImgByte bytes.Buffer
}

func NewImage(img_tpl_path string) (img *Img, err error) {
	img = new(Img)
	imgTplFile, err := os.Open(img_tpl_path)
	if err != nil {
		return
	}
	defer imgTplFile.Close()
	img.ext = path.Ext(img_tpl_path)
	switch img.ext {
	case ".jpg":
		// if img.ImgTpl, _, err = image.Decode(imgTplFile); err != nil {
		// 	return
		// }
		if img.ImgTpl, err = jpeg.Decode(imgTplFile); err != nil {
			return
		}
	case ".png":
		if img.ImgTpl, err = png.Decode(imgTplFile); err != nil {
			return
		}
	default:
		err = errors.New("暂不支持其它类型")
	}
	img.ImgOut = Image2RGBA(img.ImgTpl)
	// img.ImgOut = image.NewNRGBA(image.Rect(0, 0, img.ImgTpl.Bounds().Dx(), img.ImgTpl.Bounds().Dy()))
	// draw.Draw(img.ImgOut, img.ImgOut.Bounds(), img.ImgTpl, img.ImgTpl.Bounds().Min, draw.Over)
	return
}

//Image2RGBA Image2RGBA
func Image2RGBA(img image.Image) *image.NRGBA {
	baseSrcBounds := img.Bounds().Max
	newWidth := baseSrcBounds.X
	newHeight := baseSrcBounds.Y
	des := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight)) // 底板
	//首先将一个图片信息存入jpg
	draw.Draw(des, des.Bounds(), img, img.Bounds().Min, draw.Over)
	return des
}
func (I *Img) DrawImage(img image.Image, resize_x float32, x, y float32) (err error) {
	// 调用resize库进行图片缩放(高度填0，resize.Resize函数中会自动计算缩放图片的宽高比)
	resizeTpl := resize.Resize(uint(I.ImgTpl.Bounds().Dx()), 0, I.ImgTpl, resize.Lanczos3)
	resizeImg := resize.Resize(uint(I.ImgTpl.Bounds().Dx()*int(resize_x*10))/1000, 0, img, resize.Lanczos3) //大小
	// 将两个图片合成一张
	I.ImgOut = image.NewNRGBA(image.Rect(0, 0, resizeTpl.Bounds().Max.X, resizeTpl.Bounds().Max.Y))
	draw.Draw(I.ImgOut, I.ImgOut.Bounds(), resizeTpl, resizeTpl.Bounds().Min, draw.Over)
	draw.Draw(I.ImgOut, I.ImgOut.Bounds(), resizeImg, image.Point{X: (-1) * int(x*10) * resizeTpl.Bounds().Max.X / 1000, Y: (-1) * int(y*10) * resizeTpl.Bounds().Max.Y / 1000}, draw.Over)
	return
}

func (I *Img) DrawImageFile(img_path string, resize_x float32, x, y float32) (err error) {
	imgFile, err := os.Open(img_path)
	if err != nil {
		return
	}
	defer imgFile.Close()
	var img image.Image
	img, _, err = image.Decode(imgFile)
	// img, err = png.Decode(imgFile)
	if err != nil {
		return
	}
	err = I.DrawImage(img, resize_x, x, y)
	return
}

func (I *Img) PackImage() (err error) {
	// w := bufio.NewWriter(&I.ImgByte.)
	switch I.ext {
	case ".jpg":
		err = jpeg.Encode(&I.ImgByte, I.ImgOut, &jpeg.Options{Quality: 90})
	case ".png":
		err = png.Encode(&I.ImgByte, I.ImgOut)
	default:
		err = errors.New("未知类型")
	}
	return
}

func (I *Img) OutImage(file_path string) (err error) {
	if I.ImgOut == nil {
		err = errors.New("OutImage 不能为空")
		return
	}
	if PathExists(file_path) {
		os.Remove(file_path)
	}
	// 保存文件
	imgFile, err := os.Create(file_path)
	if err != nil {
		return
	}
	defer imgFile.Close()

	switch path.Ext(file_path) {
	case ".jpg":
		err = jpeg.Encode(imgFile, I.ImgOut, &jpeg.Options{Quality: 90})
	case ".png":
		err = png.Encode(imgFile, I.ImgOut)
	default:
		fmt.Println(I.ext)
		return
	}
	if err != nil {
		return
	}
	return
}

func NewFreeType(font_path string, size_font float64, clr color.Color) (c *freetype.Context, err error) {
	// c.SetDPI(2048)
	fontBytes, err := ioutil.ReadFile(font_path)
	if err != nil {
		return
	}
	fontFam, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return
	}
	c = freetype.NewContext()
	c.SetFont(fontFam)
	c.SetFontSize(size_font)
	// 设置文字颜色、字体、字大小
	c.SetSrc(image.NewUniform(clr))
	c.SetHinting(font.HintingNone)
	return
}
func NewFont(font_path string) (font *truetype.Font, err error) {
	fontFile, err := ioutil.ReadFile(font_path)
	if err != nil {
		return
	}
	font, err = truetype.Parse(fontFile)
	return
}
func (I *Img) DrawText(fontType *truetype.Font, size_font float64, clr color.Color, text string, x, y float32) (err error) {
	if I.ImgOut == nil {
		err = errors.New("DrawText ImgOut 不能为空")
		return
	}

	// fmt.Println("DX:", I.ImgOut.Bounds().Dx())
	// fmt.Println("SizeFont:", size_font)
	c := freetype.NewContext()
	c.SetFont(fontType)
	c.SetFontSize(size_font)
	c.SetSrc(image.NewUniform(clr))
	c.SetDPI(float64(I.ImgOut.Bounds().Dx()))
	c.SetClip(I.ImgOut.Bounds())
	c.SetDst(I.ImgOut)
	c.SetHinting(font.HintingNone)
	pt := freetype.Pt(I.ImgOut.Bounds().Dx()*int(x*10)/1000, I.ImgOut.Bounds().Dy()*int(y*10)/1000)
	// fmt.Printf("%d %d\r\n", I.ImgOut.Bounds().Dx(), I.ImgOut.Bounds().Dy())
	// fmt.Printf("%+v\r\n", pt)
	_, err = c.DrawString(text, pt)
	// fmt.Println(strings.Repeat("+", 10))
	return
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func Create() {
	img, err := NewImage("tpl_ticket.png")
	// img, err := NewImage("tpl.jpg")
	if err != nil {
		return
	}
	// err = img.DrawImageFile("9637.jpg", 23, 69.5, 81.8)
	// if err != nil {
	// 	return
	// }
	fontType, _ := NewFont("YaHeiBold.ttf")
	err = img.DrawText(fontType, 4.5, color.White, "这测试", 25, 66.5)
	if err != nil {
		return
	}
	err = img.DrawText(fontType, 3.0, color.White, "这是一个测试", 23, 71)
	if err != nil {
		return
	}
	err = img.DrawText(fontType, 4.5, color.White, "这是一个测试", 23, 76)
	if err != nil {
		return
	}
	err = img.DrawText(fontType, 0.9, color.White, "生成时间："+time.Now().Format("2006-01-02 15:04:05"), 75, 99)
	if err != nil {
		return
	}
	err = img.OutImage("out.png")
	// err = img.OutImage("out.jpg")
	if err != nil {
		return
	}
}
