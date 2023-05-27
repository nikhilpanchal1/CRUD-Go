package handlers

import (
	"os"
	"time"

	"example.com/go-fiber-api/cmd/models"
	"example.com/go-fiber-api/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Hello, ROUTEEEEEEEES!")
}

func GetItem(c *fiber.Ctx) error {
	// Logic to fetch all items from database
	items := []models.Item{}

	database.DB.Db.Find(&items)
	return c.Status(200).JSON(items)
}

// method that gets item by id
func GetItemById(c *fiber.Ctx) error {
	itemId := c.Params("ID")

	result := database.DB.Db.First(&models.Item{}, "ID = ?", itemId)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(200).JSON(result)
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
	itemId := c.Params("ID") //?NOTE: Need to change id
	//check if item exists in database
	//result := database.DB.Db.Delete(&models.Item{}, "ID = ?", c.Params(itemId)) //

	result := database.DB.Db.Delete(&models.Item{}, "ID = ?", itemId) //

	//if its not found, send status 404
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	//Assuming no errors send status 200
	return c.SendStatus(fiber.StatusOK)
}

func DeleteAll(c *fiber.Ctx) error {
	// Logic to delete all items
	database.DB.Db.Exec("DELETE FROM items")
	return c.SendString("All items deleted")
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var user models.User
	database.DB.Db.Where("username = ?", input.Username).First(&user)
	if user.ID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid login credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid login credentials",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	var jwtKey = []byte(os.Getenv("SECRET_KEY"))

	//t, err := token.SignedString([]byte(jwtKey)) //Incase we need to use hardcoded secret key
	t, err := token.SignedString(jwtKey)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// Save the token in the database
	database.DB.Db.Create(&models.Token{
		UserID: user.ID,
		Token:  t,
	})

	return c.JSON(fiber.Map{"token": t})
}

func RegisterUser(c *fiber.Ctx) error {
	type RegisterInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	input := new(RegisterInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := models.NewUser(input.Username, input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	database.DB.Db.Create(&user)

	return c.Status(200).JSON(user)
}

func GetUsers(c *fiber.Ctx) error {
	// Logic to fetch all items from database
	users := []models.User{}

	database.DB.Db.Find(&users)
	return c.Status(200).JSON(users)
}

func DeleteUser(c *fiber.Ctx) error {
	// Logic to delete an item
	userId := c.Params("ID")                                          //?NOTE: Need to change id
	result := database.DB.Db.Delete(&models.User{}, "ID = ?", userId) //

	//if its not found, send status 404
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	//Assuming no errors send status 200
	return c.SendStatus(fiber.StatusOK)
}

func GetLoggedInUsers(c *fiber.Ctx) error {
	// Query the tokens table to get all active tokens
	tokens := []models.Token{}
	database.DB.Db.Find(&tokens)

	// For each token, find the corresponding user and add them to a list
	users := []models.User{}
	for _, token := range tokens {
		user := models.User{}
		database.DB.Db.Where("ID = ?", token.UserID).First(&user)
		users = append(users, user)
	}

	return c.Status(200).JSON(users)
}

func LogoutAllUsers(c *fiber.Ctx) error {
	// Delete all tokens from the database
	database.DB.Db.Exec("DELETE FROM tokens")
	return c.SendString("All users logged out")
}
