package controller

import (
	"github.com/Vishal21121/code-room-backend-go.git/util"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var userController *mongo.Collection

func SetClient(c *mongo.Client) {
	client = c
	userController = client.Database("code-room-go").Collection("users")
}

func RegisterUser(c *fiber.Ctx) error {
	var data map[string]string
	c.BodyParser(&data)
	if data["username"] == "" || len(data["username"]) < 3 {
		return c.JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"statusCode": 400,
				"message":    "Username must be at least 3 characters",
			},
		})
	} else if data["password"] == "" || len(data["password"]) < 8 {
		return c.JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"statusCode": 400,
				"message":    "Password must be at least 8 characters",
			},
		})
	} else if data["email"] == "" {
		return c.JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"statusCode": 400,
				"message":    "Email is required",
			},
		})
	}
	var foundUser map[string]interface{}
	userController.FindOne(c.Context(), fiber.Map{"email": data["email"]}).Decode(&foundUser)
	if foundUser != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"statusCode": 400,
				"message":    "Please enter another email address",
			},
		})
	}
	hashedPassword, err := util.HashPassword(data["password"])
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"statusCode": 500,
				"message":    "Internal server error",
			},
		})
	}
	result, err := userController.InsertOne(c.Context(), fiber.Map{"username": data["username"], "password": hashedPassword, "email": data["email"], "isLoggedIn": false, "refreshToken": ""})
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"statusCode": 500,
				"message":    "Internal server error",
			},
		})
	}
	var userGot map[string]interface{}
	userController.FindOne(c.Context(), fiber.Map{"_id": result.InsertedID}).Decode(&userGot)
	// remove the password from the userGot
	delete(userGot, "password")
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"statusCode": 201,
			"message":    "User created successfully",
			"value":      userGot,
		},
	})
}
