package controller

import (
	"fmt"
	"time"

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

func LoginUser(c *fiber.Ctx) error {
	var dataReceived map[string]string
	c.BodyParser(&dataReceived)
	if dataReceived["email"] == "" {
		return c.JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"statusCode": 400,
				"message":    "Email is required",
			},
		})
	} else if dataReceived["password"] == "" {
		return c.JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"statusCode": 400,
				"message":    "Password is required",
			},
		})
	}
	fmt.Println(dataReceived["email"], dataReceived["password"])
	var foundUser map[string]string
	userController.FindOne(c.Context(), fiber.Map{"email": dataReceived["email"]}).Decode(&foundUser)
	if foundUser == nil {
		return c.JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"statusCode": 401,
				"message":    "Please enter valid credentials",
			},
		})
	}
	if !util.CheckPasswordHash(dataReceived["password"], foundUser["password"]) {
		fmt.Println("password not matched")
		return c.JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"statusCode": 401,
				"message":    "Please enter correct credentials",
			},
		})
	}
	// generate jwt token and send it to user
	accessToken, err := util.GenerateAccessToken(fiber.Map{"_id": foundUser["_id"], "username": foundUser["username"], "email": foundUser["email"]})

	if err != nil {
		fmt.Println("error in generating jwt token")
		return c.JSON(fiber.Map{
			"succes": false,
			"data": fiber.Map{
				"statusCode": 500,
				"message":    "Internal server error",
			},
		})
	}
	refreshToken, err := util.GenerateRefershToken(foundUser["_id"])
	if err != nil {
		return c.JSON(fiber.Map{
			"succes": false,
			"data": fiber.Map{
				"statusCode": 500,
				"message":    "Internal server error",
			},
		})
	}
	// update the user's refresh token in the database
	userController.UpdateOne(c.Context(), fiber.Map{"_id": foundUser["_id"]}, fiber.Map{"$set": fiber.Map{"refreshToken": refreshToken}})

	cookie := new(fiber.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = refreshToken
	cookie.MaxAge = int(15 * 24 * time.Hour.Milliseconds())
	cookie.Secure = true
	cookie.HTTPOnly = true
	cookie.SameSite = "none"
	c.Cookie(cookie)

	delete(foundUser, "password")
	delete(foundUser, "refreshToken")

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"statusCode": 200,
			"message":    "User logged in successfully",
			"value": fiber.Map{
				"accessToken":  accessToken,
				"loggedInUser": foundUser,
			},
		},
	})

}
