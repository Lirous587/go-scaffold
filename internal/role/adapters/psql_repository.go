package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"os"
	"scaffold/internal/common/orm"
	"scaffold/internal/role/model"
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

func (P *PSQLRepository) CreateRole(req *model.CreateRoleReq) (*model.Role, error) {
	ctx := context.Background()

	entity, err := orm.Roles(orm.RoleWhere.Name.EQ(req.Name)).One(ctx, P.db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if entity != nil {
		return nil, errors.New("该用角色存在")
	}

	entity = &orm.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := entity.Insert(ctx, P.db, boil.Infer()); err != nil {
		return nil, err
	}

	return &model.Role{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
	}, nil
}

func (P *PSQLRepository) DeleteRole(id int) error {
	ctx := context.Background()
	entity, err := orm.Roles(orm.RoleWhere.ID.EQ(int(id))).One(ctx, P.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("该角色不存在")
		}
		return err
	}

	_, err = entity.Delete(ctx, P.db)
	return err
}

func (P *PSQLRepository) UpdateRole(id int, req *model.UpdateRoleReq) error {
	ctx := context.Background()
	entity, err := orm.Roles(orm.RoleWhere.ID.EQ(int(id))).One(ctx, P.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("该角色不存在")
		}
		return err
	}

	entity.Name = req.Name
	entity.Description = req.Description

	_, err = entity.Update(ctx, P.db, boil.Infer())
	return err
}

func (P *PSQLRepository) AllRole() ([]model.Role, error) {
	ctx := context.Background()
	entities, err := orm.Roles().All(ctx, P.db)
	if err != nil {
		return nil, err
	}
	roles := make([]model.Role, len(entities))
	for i, entity := range entities {
		roles[i] = model.Role{
			ID:          entity.ID,
			Name:        entity.Name,
			Description: entity.Description,
		}
	}
	return roles, err
}
