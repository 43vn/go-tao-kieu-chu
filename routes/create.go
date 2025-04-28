package routes
import (
	"github.com/gofiber/fiber/v2"
	"tao-kieu-chu/api"
	"fmt"
	"html"
)
func Create(c *fiber.Ctx) error {
	text := c.FormValue("text")
	font := c.FormValue("font")
	if !api.FontExists(font) {
		errorMessage := "Font không hợp lệ"
		encodedError := api.B64E(errorMessage)
		return c.Redirect(fmt.Sprintf("/?error=%s", encodedError))
	}
	cleanString := html.UnescapeString(text)
	g := api.Gen {
		Text: cleanString,
		Font: font,
	}
	result, err := g.Encrypt()
	if err != nil {
		encodedError := api.B64E(err.Error())
		return c.Redirect(fmt.Sprintf("/?error=%s", encodedError))
	}
	return c.Redirect(fmt.Sprintf("/?hash=%s", result))
}
