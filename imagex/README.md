# imagex

图片工具包（原 `img` 包已删除并与本包合并，本包是仓库内唯一的图片工具），提供基于模板的图片合成、文字绘制以及二维码生成。

## 功能列表

### img.go（图片处理）

- `Img`：图片操作结构体（模板图 `ImgTpl`、输出图 `ImgOut`、字节缓冲 `ImgByte`）。
- `NewImage(img_tpl_path)`：从模板文件创建（支持 .jpg/.png），并转换为 NRGBA 底板。
- `Image2RGBA(img)`：将任意 image.Image 绘制到 NRGBA 底板并返回。
- `DrawImage(img, resize_x, x, y)`：将传入图片按比例缩放后叠加到模板上，`x`/`y` 为千分比坐标定位。
- `DrawImageFile(img_path, resize_x, x, y)`：从文件读取图片并调用 `DrawImage` 合成。
- `PackImage()`：将输出图编码（jpg 质量 90 / png）到内存 `ImgByte`。
- `OutImage(file_path)`：将输出图保存为文件（已存在则先删除），按扩展名选择编码。
- `NewFreeType(font_path, size_font, clr)`：读取字体文件并创建 freetype 绘制上下文。
- `NewFont(font_path)`：解析字体文件返回 `*truetype.Font`。
- `DrawText(fontType, size_font, clr, text, x, y)`：在输出图上按千分比坐标绘制文字（DPI 随图宽自适应）。
- `PathExists(path)`：判断文件路径是否存在。
- `Create()`：内部演示函数：加载模板、绘制多行文字并输出 out.png。

### qrcode.go（二维码）

- `Qrcode`：二维码结构体（字节内容 `ByteFile`、图像 `ImgQrcode`）。
- `NewQrcode(text, ext)`：生成 256x256 的二维码（容错级 M、自动选择编码模式），ext 取 `.jpg` / `.png`。
- `OutQrcode(file_path)`：将二维码图像保存到文件。

## 使用示例

```go
// 模板图上绘制文字（摘自 Create()）
img, err := imagex.NewImage("tpl_ticket.png")
if err == nil {
    fontType, _ := imagex.NewFont("YaHeiBold.ttf")
    _ = img.DrawText(fontType, 4.5, color.White, "测试文字", 25, 66.5)
    _ = img.OutImage("out.png")
}

// 生成并保存二维码
q, _ := imagex.NewQrcode("https://example.com", ".png")
_ = q.OutQrcode("qrcode.png")
```

## 备注

- 外部依赖：`github.com/golang/freetype`（文字绘制）、`github.com/nfnt/resize`（图片缩放）、`github.com/boombuler/barcode`（二维码）。
- `DrawText` 需要自备 TTF 字体文件（如 YaHeiBold.ttf）。
- 当前仅支持 jpg/png 两种格式。
