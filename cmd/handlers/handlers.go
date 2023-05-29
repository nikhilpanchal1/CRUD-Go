package handlers

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"example.com/go-fiber-api/cmd/models"
	"example.com/go-fiber-api/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Hello, ROUTEEEEEEEES!")
}

var (
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_ADDR"), // Redis server address
		Password:     "",                      // os.Getenv("DB_PASSWORD"), // redis password
		DB:           0,                       // convert string to int: strconv.Atoi(os.Getenv("DB_NUMBER")), // database number
		MinIdleConns: 4,
		PoolSize:     40,
	})
)

func GetItemById(c *fiber.Ctx) error {
	itemId := c.Params("ID")

	// Try Redis cache first
	cacheVal, err := rdb.Get(ctx, "item:"+itemId).Result()

	// If cache entry exists
	if err == nil {
		var item models.Item
		err := json.Unmarshal([]byte(cacheVal), &item)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(200).JSON(item)
	} else if err == redis.Nil {
		// Cache entry doesnt exist fetching from database
		var item models.Item
		result := database.DB.Db.First(&item, "ID = ?", itemId)
		if result.Error != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}

		// Saving to redis
		itemJson, err := json.Marshal(item) // Changed from 'result' to 'item'
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		err = rdb.Set(ctx, "item:"+itemId, itemJson, time.Hour).Err() // Cache for 1 hour, // 0 means no expiration
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(200).JSON(item) // Changed from 'result' to 'item'
	} else {
		// Redis error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
}

func GetItem(c *fiber.Ctx) error {
	items := []models.Item{}

	// Getting from redis
	cacheVal, err := rdb.Get(ctx, "items").Result()

	// If cache entry exists
	if err == nil {
		err := json.Unmarshal([]byte(cacheVal), &items)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	} else if err == redis.Nil {
		// fetching from db
		database.DB.Db.Find(&items)

		//Saving to redis
		itemsJson, err := json.Marshal(items)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		err = rdb.Set(ctx, "items", itemsJson, 0).Err() // 0 means no expiration
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	} else {
		// Redis error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(items)
}

func AddItem(c *fiber.Ctx) error {
	item := new(models.Item)
	if err := c.BodyParser(item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	database.DB.Db.Create(&item)

	// Invalidate the 'items' cache entry
	err := rdb.Del(ctx, "items").Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(item)
}

func DeleteItem(c *fiber.Ctx) error {
	// Logic to delete an item
	itemId := c.Params("ID")
	result := database.DB.Db.Delete(&models.Item{}, "ID = ?", itemId)

	//if its not found, send status 404
	if result.Error != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	// Invalidate the 'items' cache entry
	err := rdb.Del(ctx, "items").Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	//Assuming no errors send status 200
	return c.SendStatus(fiber.StatusOK)
}

func DeleteAll(c *fiber.Ctx) error {
	// Logic to delete all items
	database.DB.Db.Exec("DELETE FROM items")

	// Invalidate the 'items' cache entry
	err := rdb.Del(ctx, "items").Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

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
