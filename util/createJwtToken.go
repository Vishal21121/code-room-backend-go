package util

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateRefershToken(id string) (string, error) {
	seconds := 15 * 24 * time.Hour.Seconds()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"_id": id,
		"exp": time.Now().UTC().Add(time.Duration(seconds)).Unix(),
	})
	tokenSigned, error := token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if error != nil {
		return "", error
	}
	return tokenSigned, nil
}

func GenerateAccessToken(data fiber.Map) (string, error) {
	seconds := 15 * time.Minute.Seconds()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  time.Now().UTC().Add(time.Duration(seconds)).Unix(),
	})
	tokenSigned, error := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if error != nil {
		return "", error
	}
	return tokenSigned, nil
}
