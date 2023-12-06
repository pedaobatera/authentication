package middleware

import (
	"strings"

	"github.com/pedaobatera/monery.packages.my_authentication/util"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func ValidateTokenMiddleware(audience string, auth0Domain string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Pegar o token do cabeçalho Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Missing authorization header",
				"error":   "unauthorized",
				"code":    fiber.StatusUnauthorized,
			})
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse e validação do token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "unexpected signing method.",
					"error":   "unauthorized",
					"code":    fiber.StatusUnauthorized,
				})
			}
			cert, err := util.GetPemCert(token, auth0Domain)
			if err != nil {

				return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Error getting pem cert.",
					"error":   "unauthorized",
					"code":    fiber.StatusUnauthorized,
				})
			}
			return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token not valid.",
				"error":   "unauthorized",
				"code":    fiber.StatusUnauthorized,
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			aud := audience
			checkAud := claims.VerifyAudience(aud, false)
			if !checkAud {

				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"message": "Invalid audience.",
					"error":   "unauthorized",
					"code":    fiber.StatusUnauthorized,
				})
			}
			// Continue com o próximo middleware/handler
			c.Locals("token", token)
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token not valid.",
			"error":   "unauthorized",
			"code":    fiber.StatusUnauthorized,
		})
	}
}
