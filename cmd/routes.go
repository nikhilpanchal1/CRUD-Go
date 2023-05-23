package main

import (
	"example.com/go-fiber-api/cmd/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/items", handlers.GetItem)
	app.Post("/items", handlers.AddItem)
	app.Delete("/items", handlers.DeleteAll)
	app.Delete("/items/:id", handlers.DeleteItem) ///:id
}
