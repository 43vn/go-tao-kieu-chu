package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"

	"github.com/fogleman/gg"
)

var (
	fontSize    = 26
	padding     = 2
	extraHeight = 0
)

func GetDefault() (string, string, string) {
	return strconv.Itoa(fontSize), strconv.Itoa(padding), strconv.Itoa(extraHeight)
}

func GetFontSize(size string) int {
	if size, err := strconv.Atoi(size); err == nil {
		return size
	}
	return fontSize
}

func GetPadding(size string) int {
	if size, err := strconv.Atoi(size); err == nil {
		return size
	}
	return padding
}

func GetExtraHeight(size string) int {
	if size, err := strconv.Atoi(size); err == nil {
		return size
	}
	return extraHeight
}

func (g *Gen) CreateImageWithFont() (string, error) {
	// Tính chiều dài và chiều cao của hình ảnh
	fontSize := GetFontSize(g.FontSize)
	padding := GetPadding(g.Padding)
	extraHeight := GetExtraHeight(g.Height)
	ctx := gg.NewContext(50, 50)
	// Tính toán chiều rộng và chiều cao của văn bản
	textWidth, textHeight := ctx.MeasureString(g.Text)
	width := padding + int(textWidth*float64(padding))
	height := fontSize + int(textHeight) + extraHeight

	// Tạo lại context với kích thước vừa tính
	ctx = gg.NewContext(width, height)
	ctx.SetColor(color.White)
	ctx.Clear()
	if err := ctx.LoadFontFace(fmt.Sprintf("%s", filepath.Join(workingDir(), "fonts", g.Font)), float64(fontSize)); err != nil {
		return "", fmt.Errorf("không thể tải font")
	}
	ctx.SetColor(color.Black)
	ctx.DrawStringAnchored(g.Text, float64(width/2), float64(height/2), 0.5, 0.5)
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
