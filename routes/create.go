package routes

import (
	"fmt"
	"html"
	"tao-kieu-chu/api"

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx) error {
	text := c.FormValue("text")
	font := c.FormValue("font")
	fontSize := c.FormValue("fontSize")
	padding := c.FormValue("padding")
	height := c.FormValue("height")
	if !api.FontExists(font) {
		errorMessage := "Font không hợp lệ"
		encodedError := api.B64E(errorMessage)
		return c.Redirect(fmt.Sprintf("/?error=%s", encodedError))
	}
	cleanString := html.UnescapeString(text)
	g := api.Gen{
		Text:     cleanString,
		Font:     font,
		FontSize: fontSize,
		Padding:  padding,
		Height:   height,
	}
	result, err := g.EncodeCompressed()
	if err != nil {
		encodedError := api.B64E(err.Error())
		return c.Redirect(fmt.Sprintf("/?error=%s", encodedError))
	}
	return c.Redirect(fmt.Sprintf("/?d=%s", result))
}
