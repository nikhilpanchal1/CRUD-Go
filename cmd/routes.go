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
	//app.Get("/items/:id", handlers.GetItemById)
	app.Delete("/items/:id", handlers.DeleteItem) ///:id
	app.Post("/login", handlers.Login)
	app.Get("/user", handlers.GetLoggedInUsers) //Ge all logged in users
	app.Post("/users", handlers.RegisterUser)
	app.Get("/users", handlers.GetUsers)
	app.Delete("/users/:id", handlers.DeleteUser)
	app.Post("/users/logout", handlers.LogoutAllUsers)
}
