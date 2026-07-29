package controllers

import (
	"los_andes/database"
	"los_andes/models"

	"github.com/gofiber/fiber/v2"
)

func ObtenerDashboard(c *fiber.Ctx) error {
	var (
		stats models.DashboardResponse
		conn  = database.GetDB()
	)

	query := `
		SELECT 
			COALESCE((SELECT COUNT(*) FROM equipos), 0) AS equipos,
			COALESCE((SELECT COUNT(*) FROM entregas), 0) AS entregas,
			COALESCE((SELECT COUNT(*) FROM clientes), 0) AS clientes,
			COALESCE((SELECT COUNT(*) FROM log_error), 0) AS errores`

	err := conn.QueryRow(query).Scan(
		&stats.Equipos,
		&stats.Entregas,
		&stats.Clientes,
		&stats.Errores,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al consultar las estadísticas del dashboard: " + err.Error(),
		})
	}

	return c.JSON(stats)
}
