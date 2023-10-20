package controllers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/osmansam/react/configs"
	"github.com/osmansam/react/models"
	"github.com/osmansam/react/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var tokenCollection *mongo.Collection = configs.GetCollection(configs.DB, "tokens")

func Register(c *fiber.Ctx) error {
	var user models.User

	// Parse the request body into the user model
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Check if email exists
	count, _ := userCollection.CountDocuments(context.Background(), bson.M{"email": user.Email})
	if count > 0 {
		return c.Status(http.StatusBadRequest).SendString("Email already exists")
	}

	// First registered user is an admin, otherwise a user
	// Here you need to decide how to implement roles in your Go model. I'm using a hypothetical "Role" field.
	totalUsers, _ := userCollection.CountDocuments(context.Background(), bson.M{})
	user.Role = "user"
	if totalUsers == 0 {
		user.Role = "admin"
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error hashing password")
	}
	user.Password = hashedPassword

	// Insert user into database
	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Error creating user")
	}

	return c.Status(http.StatusCreated).SendString("Success! Account created")
}

