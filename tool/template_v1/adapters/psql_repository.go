package adapters

import (
	"database/sql"
	"fmt"
	"os"
	"scaffold/internal/user/model"
)

type PSQLRepository struct {
	db *sql.DB
}

func NewPSQLRepository() *PSQLRepository {
	host := os.Getenv("PSQL_HOST")
	port := os.Getenv("PSQL_PORT")
	user := os.Getenv("PSQL_USERNAME")
	password := os.Getenv("PSQL_PASSWORD")
	dbname := os.Getenv("PSQL_DB_NAME")
	sslmode := os.Getenv("PSQL_SSL_MODE")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return &PSQLRepository{
		db: db,
	}
}

func (P PSQLRepository) FindByEmail(email string) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (P PSQLRepository) FindByGithubID(id int) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (P PSQLRepository) Create(user *model.User) error {
	//TODO implement me
	panic("implement me")
}
