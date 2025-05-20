package adapters

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type PSQLRepository struct {
	db *sql.DB
}

func NewPSQLRepository() PSQLRepository {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := "require"

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return PSQLRepository{
		db: db,
	}
}
