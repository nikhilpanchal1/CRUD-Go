package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(" Hello, Not!")
	})
	app.Listen(":3000")
}

/*
import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	//import model.go
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/items", getItems)
	app.Post("/items", createItem)
	app.Delete("/items/:id", deleteItem)
	app.Post("/login", loginUser)
	app.Get("/user", getUser)

	app.Listen(":3000")
}
*/
