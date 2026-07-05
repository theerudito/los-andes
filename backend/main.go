package main

import (
	"fmt"
	"log"
	"los_andes/database"
	"los_andes/router"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No se pudo cargar el archivo .env, usando variables del sistema")
	}

	// 2. Inicializar base de datos (Ejecuta DDL y DML si es necesario)
	database.InitDB()

	// Obtener la instancia de la base de datos para asegurar el cierre al apagar el servidor
	db := database.GetDB()
	defer func() {
		log.Println("💾 Cerrando conexión de la base de datos...")
		db.Close()
	}()

	// 3. Crear aplicación Fiber
	app := fiber.New(fiber.Config{
		AppName: "Gestión de Mantenimiento Los Andes v1.0",
	})

	// 4. Configurar rutas de la API
	router.SetupRoutes(app)

	// 5. Determinar puerto de escucha
	port := os.Getenv("PortServer")
	if port == "" {
		port = "3000" // Puerto por defecto si no está en el .env
	}

	// 6. Encender el servidor
	log.Printf("🚀 Servidor corriendo en el puerto %s", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("❌ Error al levantar el servidor: %v", err)
	}
}
