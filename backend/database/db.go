package database

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "modernc.org/sqlite"
)

type MyDB struct {
	DB *sql.DB
}

var (
	instance *MyDB
	once     sync.Once
)

func InitDB() {

	once.Do(func() {

		db, err := sql.Open("sqlite", "istla.db")
		if err != nil {
			log.Fatalf("❌ Error al abrir SQLite: %v", err)
		}

		if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
			log.Fatalf("❌ Error activando foreign keys: %v", err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("❌ No se pudo conectar a SQLite: %v", err)
		}

		log.Println("✅ Conectado a SQLite exitosamente.")

		ejecutarScriptSQL(db, "ddl.sql")

		if debeEjecutarDML(db) {
			log.Println("🔹 Detectada base de datos nueva. Ejecutando dml.sql...")
			ejecutarScriptSQL(db, "dml.sql")
			log.Println("✅ dml.sql aplicado correctamente.")
		} else {
			log.Println("skip ⏭️ dml.sql omitido: Los datos iniciales ya existen en la base de datos.")
		}

		instance = &MyDB{DB: db}
	})
}

func GetDB() *sql.DB {
	if instance == nil {
		InitDB()
	}
	return instance.DB
}

func debeEjecutarDML(db *sql.DB) bool {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM config_inicial WHERE inicializado = 1").Scan(&count)
	if err != nil {
		return true
	}
	return count == 0
}

func ejecutarScriptSQL(db *sql.DB, rutaArchivo string) {
	contenido, err := os.ReadFile(rutaArchivo)
	if err != nil {
		log.Fatalf("❌ Error crítico al leer el archivo %s: %v", rutaArchivo, err)
	}

	_, err = db.Exec(string(contenido))
	if err != nil {
		log.Fatalf("❌ Error crítico al ejecutar las sentencias de %s: %v", rutaArchivo, err)
	}
}
