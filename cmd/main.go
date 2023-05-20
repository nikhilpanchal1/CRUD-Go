package main

import (
	"example.com/go-fiber-api/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	setupRoutes(app)

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
