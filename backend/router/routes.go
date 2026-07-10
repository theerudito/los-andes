package router

import (
	"los_andes/controllers"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App) {

	allowedOrigins := map[string]bool{
		os.Getenv("URL_Frontend"): true,
		"":                        true,
	}

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowOriginsFunc: func(origin string) bool {
			return allowedOrigins[origin]
		},
	}))

	api := app.Group("/api")

	v1 := api.Group("/v1")

	// clientes
	v1.Get("/cliente", controllers.ObtenerClientes)
	v1.Get("/cliente/:id", controllers.ObtenerCliente)
	v1.Get("/cliente-dni/:identificacion", controllers.ObtenerClientePorIdentificacion)
	v1.Post("/cliente", controllers.CrearCliente)
	v1.Put("/cliente", controllers.ModificarCliente)
	v1.Delete("/cliente/:id", controllers.EliminarCliente)

	// usuarios
	v1.Get("/usuario", controllers.ObtenerUsuarios)
	v1.Get("/usuario/:id", controllers.ObtenerUsuario)
	v1.Get("/usuario-dni/:identificacion", controllers.ObtenerUsuarioPorIdentificacion)
	v1.Post("/usuario", controllers.CrearUsuario)
	v1.Put("/usuario", controllers.ModificarUsuario)
	v1.Delete("/usuario/:id", controllers.EliminarUsuario)

	// marca
	v1.Get("/marca", controllers.ObtenerMarcas)
	v1.Get("/marca/:id", controllers.ObtenerMarca)
	v1.Post("/marca", controllers.CrearMarca)
	v1.Put("/marca", controllers.ModificarMarca)
	v1.Delete("/marca/:id", controllers.EliminarMarca)

	// estados
	v1.Get("/estado", controllers.ObtenerEstados)
	v1.Get("/estado/:id", controllers.ObtenerEstado)

	// equipos
	v1.Get("/equipo", controllers.ObtenerEquipos)
	v1.Get("/equipo/:id", controllers.ObtenerEquipo)
	v1.Get("/equipo/:tipo/:valor", controllers.ObtenerEquipoPorTipo)
	v1.Post("/equipo", controllers.CrearEquipo)
	v1.Put("/equipo", controllers.ModificarEquipo)
	v1.Delete("/equipo/:id", controllers.EliminarEquipo)

	// historial
	v1.Get("/historial/:id", controllers.ConsultarHistorialEquipo)
	v1.Put("/historial", controllers.ActualizarEstadoEquipo)

	// pagos
	v1.Get("/pago/:id", controllers.ConsultarCuentaEquipo)
	v1.Put("/pago", controllers.ActualizarCuentaEquipo)

	// entregas
	v1.Get("/entrega/:id", controllers.ConsultarEntregaPorEquipo)
	v1.Post("/entrega", controllers.RegistrarEntrega)

}
