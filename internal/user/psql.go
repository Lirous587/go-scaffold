package user

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type PSQL struct {
	db *sql.DB
}

func NewPsql() *PSQL {
	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSL_MODE")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return &PSQL{
		db: db,
	}
}
