package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type MyDB struct {
	DB *sql.DB
}

var instance *MyDB

func InitDB() {
	if instance != nil {
		return
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("ServerDB"), os.Getenv("PortDB"), os.Getenv("UserDB"), os.Getenv("PasswordBD"), os.Getenv("NameDB"))

	db, err := sql.Open(os.Getenv("DriverDB"), dsn)

	if err != nil {
		log.Fatalf("Error al abrir la base de datos: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	log.Println("✅ Conectado a la base de datos:", os.Getenv("DriverDB"))

	instance = &MyDB{DB: db}
}

func GetDB() *sql.DB {
	if instance == nil {
		InitDB()
	}
	return instance.DB
}
