package routes

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUpUserRoutes(app *fiber.App, client *mongo.Client) {
	userRoutes := app.Group("/api/v1/users")
}
