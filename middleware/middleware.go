package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("token")
	if cookie == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	jwtSecret := os.Getenv("jwtSecret")
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}
