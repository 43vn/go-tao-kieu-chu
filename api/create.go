package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"path/filepath"

	"github.com/fogleman/gg"
)

func CreateImageWithFont(text, font string) (string, error) {
	// Tạo hình ảnh
	const fontSize = 26
	padding := 2
	// Tính chiều dài và chiều cao của hình ảnh
	ctx := gg.NewContext(50, 50)
	// Tính toán chiều rộng và chiều cao của văn bản
	textWidth, textHeight := ctx.MeasureString(text)
	width := padding + int(textWidth*float64(padding))
	height := fontSize + int(textHeight)

	// Tạo lại context với kích thước vừa tính
	ctx = gg.NewContext(width, height)
	ctx.SetColor(color.White)
	ctx.Clear()
	if err := ctx.LoadFontFace(fmt.Sprintf("%s", filepath.Join(workingDir(), "fonts", font)), fontSize); err != nil {
		return "", fmt.Errorf("không thể tải font")
	}
	ctx.SetColor(color.Black)
	ctx.DrawStringAnchored(text, float64(width/2), float64(height/2), 0.5, 0.5)
	// Tạo buffer để lưu hình ảnh
	var buf bytes.Buffer
	// Đảm bảo giải phóng bộ nhớ ngay sau khi sử dụng
	defer buf.Reset()
	if err := ctx.EncodePNG(&buf); err != nil {
		return "", fmt.Errorf("không thể mã hóa hình ảnh thành PNG: %v", err)
	}

	// Chuyển đổi hình ảnh thành chuỗi base64
	base64String := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Trả về chuỗi base64 URL
	return "data:image/png;base64," + base64String, nil
}
