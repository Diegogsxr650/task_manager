package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // Importa el driver de MySQL
)

var DB *sql.DB

func Connect() {
	var err error
	dsn := "root:Villabalter1@tcp(127.0.0.1:3306)/task_manager?parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
	}
}
