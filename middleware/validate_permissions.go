package middleware

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func CheckPermissionsMiddleware(requiredPermissions string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token, ok := c.Locals("token").(*jwt.Token)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Token missing or not valid."})
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		permissions, ok := claims["permissions"].([]interface{})
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "You don't have the required permission."})
		}

		context, ok := claims["monery-manager/context"].(interface{}).(map[string]interface{})

		for _, permission := range permissions {
			if permission == requiredPermissions {
			}
			c.Locals("context", context)
			c.Locals("claims", claims)
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "You don't have the required permission."})
	}
}
