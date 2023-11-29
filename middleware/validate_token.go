package middleware

import (
	"authentication/util"
	"errors"
	"net/http"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func ValidateTokenMiddleware(audience string, auth0Domain string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Pegar o token do cabeçalho Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {

			return c.Status(http.StatusUnauthorized).SendString("Missing authorization header")
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse e validação do token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}
			cert, err := util.GetPemCert(token, auth0Domain)
			if err != nil {

				return nil, c.Status(http.StatusUnauthorized).JSON(err)

			}
			return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		})
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(err)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			aud := audience
			checkAud := claims.VerifyAudience(aud, false)
			if !checkAud {

				return c.Status(http.StatusUnauthorized).JSON(err)
			}
			// Continue com o próximo middleware/handler
			c.Locals("token", token)
			return c.Next()
		}

		return c.Status(http.StatusUnauthorized).SendString("Invalid token")
	}
}
