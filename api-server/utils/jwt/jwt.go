package jwt

import (
	"time"

	"github.com/ahsansaif47/blockchain-address-watcher/api-server/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(config.GetConfig().JWTSecret)

type Claims struct {
	Email string
	jwt.RegisteredClaims
}

func GenerateJWT(email string) (string, error) {
	expTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "home-kitchens",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// I wont be needing this in the auth service but this will be used in other services
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if !token.Valid {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Locals("email", claims.Email)

		return c.Next()
	}
}
