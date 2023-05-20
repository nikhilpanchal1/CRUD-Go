package main

import (
	"example.com/go-fiber-api/cmd/handlers"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/items", handlers.GetItems)
	app.Get("/item", handlers.GetItem)
	app.Post("/item", handlers.AddItem)
	app.Post("deleteAll", handlers.DeleteAll)
}
