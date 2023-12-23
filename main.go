package main

import (
	"log"

	"github.com/Vishal21121/code-room-backend-go.git/db"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	Client := db.Init()
	app.Get("/", func(c *fiber.Ctx) error {
		c.SendStatus(200)
		return c.JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"statusCode": 200,
				"value":      "Everything working fine",
			},
		})
	})

	log.Printf("\nServer is listening on port %s", ":8080")
	log.Fatal(app.Listen(":8080"))
}
