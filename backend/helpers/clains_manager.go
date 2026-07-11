package helpers

import (
	"los_andes/models"

	"github.com/gofiber/fiber/v2"
)

func ReadClaims(c *fiber.Ctx) (*models.CustomClaims, error) {

	user := c.Locals("user")

	claims, ok := user.(*models.CustomClaims)

	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Usuario no autenticado o claims inválidos")
	}

	return claims, nil
}
