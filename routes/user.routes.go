package routes

import (
	"github.com/Vishal21121/code-room-backend-go.git/controller"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpUserRoutes(app *fiber.App, client *mongo.Client) {
	controller.SetClient(client)
	userRoutes := app.Group("/api/v1/users")
	userRoutes.Post("/register", controller.RegisterUser)
	userRoutes.Post("/login", controller.LoginUser)
}
