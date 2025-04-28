package routes
import (
	"github.com/gofiber/fiber/v2"
	"tao-kieu-chu/api"
)
func Home(c *fiber.Ctx) error {
	fonts := api.SearchFonts()
	var errorMsg, imageURL, text, font string
	errB64 := c.Query("error")
	if errB64 != "" {
		errorMsg = api.B64D(errB64)
	}else{
		hash := c.Query("hash")
		if hash != "" {
			g, err := api.DecryptGenHex(hash)
			if err == nil {
				text = g.Text
				font = g.Font
				url, err := api.CreateImageWithFont(text, font)
				if err == nil && url != ""{
					imageURL = url
				}else{
					errorMsg = err.Error()
				}
			}
		}
	}
	return c.Render("home", fiber.Map{
				"Title": "Home",
				"Fonts": fonts,
				"errorMsg": errorMsg,
				"imageURL": imageURL,
				"text": text,
				"font": font,
			}, "layout")
}
