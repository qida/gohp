package img

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type Qrcode struct {
	ByteFile  []byte
	ImgQrcode image.Image
	ext       string
}

//生成二维码
func NewQrcode(text string, ext string) (qrCode Qrcode, err error) {
	qrCode.ext = ext
	img, _ := qr.Encode(text, qr.M, qr.Auto)
	img, _ = barcode.Scale(img, 256, 256)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	switch qrCode.ext {
	case ".jpg":
		err = jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
	case ".png":
		err = png.Encode(w, img)
	default:
	}
	qrCode.ImgQrcode = img
	qrCode.ByteFile = b.Bytes()
	return
}

//保存二维码
func (Q *Qrcode) OutQrcode(file_path string) (err error) {
	imgFile, _ := os.Create(file_path)
	defer imgFile.Close()
	switch Q.ext {
	case ".jpg":
		err = jpeg.Encode(imgFile, Q.ImgQrcode, &jpeg.Options{Quality: 90})
	case ".png":
		err = png.Encode(imgFile, Q.ImgQrcode)
	default:
		fmt.Println(Q.ext)
	}
	return
}
