package adapters

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"os"
	"scaffold/internal/common/orm"
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
	ctx := context.Background()
	user, err := orm.Users(orm.UserWhere.Email.EQ(email)).One(ctx, P.db)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		GithubID: user.GithubID,
	}, err
}

func (P PSQLRepository) FindByGithubID(id string) (*model.User, error) {
	ctx := context.Background()
	user, err := orm.Users(orm.UserWhere.GithubID.EQ(id)).One(ctx, P.db)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		GithubID: user.GithubID,
	}, err
}

func (P PSQLRepository) Register(u *model.User) (*model.User, error) {
	ctx := context.Background()
	user := &orm.User{
		Name:     u.Name,
		GithubID: u.GithubID,
		Email:    u.Email,
	}
	if err := user.Insert(ctx, P.db, boil.Infer()); err != nil {
		return nil, err
	}

	return &model.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		GithubID: user.GithubID,
	}, nil
}
