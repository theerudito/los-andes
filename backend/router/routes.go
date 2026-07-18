package router

import (
	"los_andes/controllers"
	"los_andes/helpers"
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

	v1.Post("/usuario/login", controllers.LoginUsuario)
	v1.Put("/usuario/reset", controllers.ResetUsuario)

	// usuarios
	protectedUsuarios := v1.Group("/usuario", helpers.JWTMiddleware())
	protectedUsuarios.Get("/", controllers.ObtenerUsuarios)
	protectedUsuarios.Get("/:id", controllers.ObtenerUsuario)
	protectedUsuarios.Get("/dni/:identificacion", controllers.ObtenerUsuarioPorIdentificacion)
	protectedUsuarios.Post("/", controllers.CrearUsuario)
	protectedUsuarios.Put("/", controllers.ModificarUsuario)
	protectedUsuarios.Delete("/:id", controllers.EliminarUsuario)

	// clientes
	protectedClientes := v1.Group("/cliente", helpers.JWTMiddleware())
	protectedClientes.Get("/", controllers.ObtenerClientes)
	protectedClientes.Get("/:id", controllers.ObtenerCliente)
	protectedClientes.Get("/dni/:identificacion", controllers.ObtenerClientePorIdentificacion)
	protectedClientes.Post("/", controllers.CrearCliente)
	protectedClientes.Put("/", controllers.ModificarCliente)
	protectedClientes.Delete("/:id", controllers.EliminarCliente)
	protectedClientes.Post("/reportes", controllers.ReporteCliente)

	// marca
	protectedMarcas := v1.Group("/marca", helpers.JWTMiddleware())
	protectedMarcas.Get("/", controllers.ObtenerMarcas)
	protectedMarcas.Get("/:id", controllers.ObtenerMarca)
	protectedMarcas.Post("/", controllers.CrearMarca)
	protectedMarcas.Put("/", controllers.ModificarMarca)
	protectedMarcas.Delete("/:id", controllers.EliminarMarca)

	// estados
	protectedEstados := v1.Group("/estado", helpers.JWTMiddleware())
	protectedEstados.Get("/", controllers.ObtenerEstados)
	protectedEstados.Get("/:id", controllers.ObtenerEstado)

	// equipos
	protectedEquipos := v1.Group("/equipo", helpers.JWTMiddleware())
	protectedEquipos.Get("/", controllers.ObtenerEquipos)
	protectedEquipos.Get("/:id", controllers.ObtenerEquipo)
	protectedEquipos.Get("/:tipo/:valor", controllers.ObtenerEquipoPorTipo)
	protectedEquipos.Post("/", controllers.CrearEquipo)
	protectedEquipos.Put("/", controllers.ModificarEquipo)
	protectedEquipos.Delete("/:id", controllers.EliminarEquipo)
	protectedEquipos.Post("/reportes", controllers.ReporteEquipos)

	// historial
	protectedHistorial := v1.Group("/historial", helpers.JWTMiddleware())
	protectedHistorial.Get("/:id", controllers.ConsultarHistorialEquipo)
	protectedHistorial.Put("/", controllers.ActualizarEstadoEquipo)

	// pagos
	protectedPagos := v1.Group("/pago", helpers.JWTMiddleware())
	protectedPagos.Get("/:id", controllers.ConsultarCuentaEquipo)
	protectedPagos.Put("/actualizar", controllers.ActualizarCuentaEquipo)
	protectedPagos.Post("/procesar", controllers.ProcesarEntregaEquipo)

	// entregas
	protectedEntregas := v1.Group("/entrega", helpers.JWTMiddleware())
	protectedEntregas.Get("/:id", controllers.ConsultarEntregaPorEquipo)
	protectedEntregas.Post("/", controllers.RegistrarEntrega)

	// logs error
	protectedLogsError := v1.Group("/logs-error", helpers.JWTMiddleware())
	protectedLogsError.Post("/", controllers.ObtenerLogsError)

	// logs ok
	protectedLogsOK := v1.Group("/logs-ok", helpers.JWTMiddleware())
	protectedLogsOK.Post("/", controllers.ObtenerLogsOk)

}
