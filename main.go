package main

import (
	"fmt"
	"net/http"
	"os"
	"tao-kieu-chu/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/handlebars/v2"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	engine := handlebars.New("./views", ".hbs")
	engine.Reload(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:         http.Dir("./assets"),
		Browse:       true,
		Index:        "index.html",
		NotFoundFile: "404.html",
		MaxAge:       3600,
	}))
	app.Get("/", routes.Home)
	app.Post("/create", routes.Create)
	app.Listen(fmt.Sprintf(":%s", port))
}
