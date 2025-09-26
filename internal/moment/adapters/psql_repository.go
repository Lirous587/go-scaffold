package adapters

import (
	"blog-v4/internal/common/orm"
	"blog-v4/internal/common/reskit/codes"
	"blog-v4/internal/common/utils"
	"blog-v4/internal/moment/domain"
	"database/sql"
	"errors"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type PSQLMomentRepository struct {
}

func NewPSQLMomentRepository() domain.MomentRepository {
	return &PSQLMomentRepository{}
}

func (repo *PSQLMomentRepository) FindByID(id int64) (*domain.Moment, error) {
	moment, err := orm.FindMomentG(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrMomentNotFound
		}
		return nil, err
	}
	return ormMomentToDomain(moment), err
}

func (repo *PSQLMomentRepository) Create(moment *domain.Moment) (*domain.Moment, error) {
	ormMoment := domainMomentToORM(moment)

	if err := ormMoment.InsertG(boil.Infer()); err != nil {
		return nil, err
	}

	return ormMomentToDomain(ormMoment), nil
}

func (repo *PSQLMomentRepository) Delete(id int64) error {
	rows, err := orm.Moments(qm.Where("id = ?", id)).DeleteAllG(false)

	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrMomentNotFound
	}
	return nil
}

func (repo *PSQLMomentRepository) Update(moment *domain.Moment) (*domain.Moment, error) {
	ormMoment := domainMomentToORM(moment)

	rows, err := ormMoment.UpdateG(boil.Infer())

	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, codes.ErrMomentNotFound
	}

	return ormMomentToDomain(ormMoment), nil
}

func (repo *PSQLMomentRepository) ListMoments(query *domain.MomentQuery) (*domain.MomentPages, error) {
	var whereMods []qm.QueryMod
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		whereMods = append(whereMods, qm.Where("(title LIKE ? OR content LIKE ? OR location LIKE ?)", like, like, like))
	}
	// 1.计算count
	count, err := orm.Moments(whereMods...).CountG()
	if err != nil {
		return nil, err
	}
	totalPages, err := utils.ComputePages(query.Page, query.PageSize, count)
	if err != nil {
		return nil, err
	}

	// 2.计算offset
	offset, err := utils.ComputeOffset(query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	pageMods := append(whereMods, qm.Offset(offset), qm.Limit(query.PageSize))

	moment, err := orm.Moments(pageMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.MomentPages{
		Pages: totalPages,
		List:  ormMomentsToDomain(moment),
	}, nil
}

func (repo *PSQLMomentRepository) RandomN(count int8) ([]*domain.Moment, error) {
	mods := []qm.QueryMod{qm.Limit(int(count)), qm.OrderBy("RANDOM()")}
	slice, err := orm.Moments(mods...).AllG()
	if err != nil {
		return nil, err
	}
	return ormMomentsToDomain(slice), nil
}
