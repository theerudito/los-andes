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

	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env")
	}

	app := fiber.New()

	database.InitDB()

	defer database.GetDB().Close()

	router.SetupRoutes(app)

	_ = app.Listen(fmt.Sprintf(":%s", os.Getenv("PortServer")))

}
