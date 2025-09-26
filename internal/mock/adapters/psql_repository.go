package adapters

import (
	"database/sql"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/pkg/errors"
	"scaffold/internal/common/reskit/codes"
	"scaffold/internal/mock/domain"
	"scaffold/internal/common/orm"
	"scaffold/internal/common/utils"
)

type PSQLMockRepository struct {
}

func NewPSQLMockRepository() domain.MockRepository {
	return &PSQLMockRepository{}
}

func (repo *PSQLMockRepository) FindByID(id int64) (*domain.Mock, error) {
	ormMock, err := orm.FindMockG(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrMockNotFound
		}
		return nil, err
	}
	return ormMockToDomain(ormMock), nil
}

func (repo *PSQLMockRepository) Create(mock *domain.Mock) (*domain.Mock,error)  {
	ormMock := domainMockToORM(mock)

	if err := ormMock.InsertG(boil.Infer()); err != nil {
		return nil, err
	}

	return ormMockToDomain(ormMock), nil
}

func (repo *PSQLMockRepository) Update(mock *domain.Mock) (*domain.Mock,error) {
	ormMock := domainMockToORM(mock)

	rows, err := ormMock.UpdateG(boil.Infer())

	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, codes.ErrMockNotFound
	}

	return ormMockToDomain(ormMock), nil
}

func (repo *PSQLMockRepository) Delete(id int64) error {
	rows, err := orm.Mocks(qm.Where("id = ?", id)).DeleteAllG(false)

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrMockNotFound
	}
	return nil
}

func (repo *PSQLMockRepository) List(query *domain.MockQuery) (*domain.MockList, error) {
	var whereMods []qm.QueryMod
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		whereMods = append(whereMods, qm.Where("(title LIKE ? OR description LIKE ?)", like, like))
	}
	// 1.计算total
	total, err := orm.Mocks(whereMods...).CountG()
	if err != nil {
		return nil, err
	}

	// 2.计算offset
	offset, err := utils.ComputeOffset(query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	listMods := append(whereMods, qm.Offset(offset), qm.Limit(query.PageSize))

	// 3.查询数据
	mock, err := orm.Mocks(listMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.MockList{
		Total: total,
		List:  ormMocksToDomain(mock),
	}, nil
}
