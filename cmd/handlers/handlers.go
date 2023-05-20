package handlers

import (
	"example.com/go-fiber-api/cmd/models"
	"example.com/go-fiber-api/database"
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Hello, ROUTEEEEEEEES!")
}

func GetItems(c *fiber.Ctx) error {
	// Logic to fetch all items from database
	items := []models.Item{}

	database.DB.Db.Find(&items)
	return c.Status(200).JSON(items)
}

func GetItem(c *fiber.Ctx) error {
	// Logic to fetch items from database
	return nil
}
func AddItem(c *fiber.Ctx) error {
	// Logic to add a new item
	item := new(models.Item)
	if err := c.BodyParser(item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	database.DB.Db.Create(&item)

	return c.Status(200).JSON(item)
}

func DeleteItem(c *fiber.Ctx) error {
	// Logic to delete an item
	return nil
}

func Login(c *fiber.Ctx) error {
	// Logic to log in the user
	return nil
}

func GetLoggedUser(c *fiber.Ctx) error {
	// Logic to get the logged in user
	return nil
}

func DeleteAll(c *fiber.Ctx) error {
	// Logic to delete all items
	database.DB.Db.Exec("DELETE FROM items")
	return c.SendString("All items deleted")
}
