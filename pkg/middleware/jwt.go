package middleware

import (
    "os"
    jwtware "github.com/gofiber/jwt/v3"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

func JWTProtected(role string) fiber.Handler {
    return jwtware.New(jwtware.Config{
        SigningKey:   []byte(os.Getenv("JWT_SECRET")),
        ErrorHandler: jwtError,
        SuccessHandler: func(c *fiber.Ctx) error {
            user := c.Locals("user")
            token, ok := user.(*jwt.Token)
            if !ok {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
            }
            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid claims"})
            }
            if role != "" && claims["role"] != role {
                return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden role"})
            }
            return c.Next()
        },
    })
}

func jwtError(c *fiber.Ctx, err error) error {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
        "error": "unauthorized",
    })
}