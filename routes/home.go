package routes

import (
	"tao-kieu-chu/api"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	fonts := api.SearchFonts()
	var errorMsg, imageURL, text, font string
	fontSize, padding, height := api.GetDefault()
	errB64 := c.Query("error")
	if errB64 != "" {
		errorMsg = api.B64D(errB64)
	} else {
		hash := c.Query("d")
		if hash != "" {
			g, err := api.DecodeCompressed(hash)
			if err == nil {
				text = g.Text
				font = g.Font
				fontSize = g.FontSize
				padding = g.Padding
				height = g.Height
				url, err := g.CreateImageWithFont()
				if err == nil && url != "" {
					imageURL = url
				} else {
					errorMsg = err.Error()
				}
			}
		}
	}

	return c.Render("home", fiber.Map{
		"Title":    "Home",
		"Fonts":    fonts,
		"errorMsg": errorMsg,
		"imageURL": imageURL,
		"text":     text,
		"font":     font,
		"fontSize": fontSize,
		"padding":  padding,
		"height":   height,
	}, "layout")
}
