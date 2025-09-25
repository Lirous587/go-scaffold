package adapters

import (
	"scaffold/internal/common/orm"
	"scaffold/internal/common/reskit/codes"
	"scaffold/internal/common/utils"
	"scaffold/internal/img/domain"
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"time"
)

type PSQLImgRepository struct {
}

func NewPSQLImgRepository() domain.ImgRepository {
	return &PSQLImgRepository{}
}

func (repo *PSQLImgRepository) getExecutor(tx ...*sql.Tx) boil.Executor {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0]
	}
	return boil.GetDB()
}

func (repo *PSQLImgRepository) FindByID(id int64, deleted ...bool) (*domain.Img, error) {
	selectDeleted := false

	if len(deleted) > 0 && deleted[0] == true {
		selectDeleted = true
	}

	var whereMods []qm.QueryMod
	whereMods = append(whereMods, qm.Where("id = ?", id))

	if selectDeleted {
		whereMods = append(whereMods, qm.WithDeleted())
	}

	ormImg, err := orm.Imgs(whereMods...).OneG()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrImgNotFound
		}
		return nil, err
	}

	return ormImgToDomain(ormImg, selectDeleted), err
}

func (repo *PSQLImgRepository) ExistByPath(path string) (bool, error) {
	exist, err := orm.Imgs(qm.Where("path = ?", path)).ExistsG()
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (repo *PSQLImgRepository) WithTX(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if opts == nil {
		opts = &sql.TxOptions{
			Isolation:	sql.LevelReadCommitted,
			ReadOnly:	false,
		}
	}
	return boil.BeginTx(ctx, opts)
}

func (repo *PSQLImgRepository) Create(img *domain.Img, categoryID int64) (*domain.Img, error) {
	ormImg := domainImgToORM(img)

	var category *domain.Category
	var err error
	// 如果有分类ID，在插入前设置
	if categoryID != 0 {
		ormImg.CategoryID.Valid = true
		ormImg.CategoryID.Int64 = categoryID

		category, err = repo.FindCategoryByID(categoryID)
		if err != nil {
			return nil, err
		}
	}

	if category != nil {
		ormImg.Path = category.Prefix + "/" + ormImg.Path
	}

	if err := ormImg.InsertG(boil.Infer()); err != nil {
		return nil, err
	}

	return ormImgToDomain(ormImg), nil
}

func (repo *PSQLImgRepository) Delete(id int64, hard bool) error {
	rows, err := orm.Imgs(qm.Where("id = ?", id), qm.WithDeleted()).DeleteAllG(hard)
	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrImgNotFound
	}
	return nil
}

func (repo *PSQLImgRepository) Restore(id int64) (*domain.Img, error) {
	rows, err := orm.Imgs(qm.WithDeleted(), qm.Where("id = ?", id)).UpdateAllG(orm.M{
		orm.ImgColumns.DeletedAt:	nil,
		orm.ImgColumns.UpdatedAt:	time.Now(),
	})
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, codes.ErrImgNotFound
	}

	img, err := repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (repo *PSQLImgRepository) OrderByIDDesc() qm.QueryMod {
	return qm.OrderBy(orm.ImgColumns.ID + " DESC")
}

func (repo *PSQLImgRepository) OrderByUpdatedAtDesc() qm.QueryMod {
	return qm.OrderBy(orm.ImgColumns.UpdatedAt + " DESC")
}

func (repo *PSQLImgRepository) List(query *domain.ImgQuery) (*domain.ImgPages, error) {
	var whereMods []qm.QueryMod
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		whereMods = append(whereMods, qm.Where("description ILIKE ?", like))
	}

	if query.CategoryID != 0 {
		whereMods = append(whereMods, qm.Where(fmt.Sprintf("%s = ?", orm.ImgColumns.CategoryID), query.CategoryID))
	}

	if query.Deleted == true {
		whereMods = append(whereMods, qm.WithDeleted())
		whereMods = append(whereMods, qm.Where("deleted_at is not null"))
	}

	// 1.计算count
	count, err := orm.Imgs(whereMods...).CountG()
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

	pageMods := append(whereMods, qm.Offset(offset), qm.Limit(query.PageSize), repo.OrderByUpdatedAtDesc(), repo.OrderByIDDesc())

	imgs, err := orm.Imgs(pageMods...).AllG()
	if err != nil {
		return nil, err
	}

	return &domain.ImgPages{
		Pages:	totalPages,
		List:	ormImgsToDomain(imgs, query.Deleted),
	}, nil
}

func (repo *PSQLImgRepository) CreateCategory(category *domain.Category) (*domain.Category, error) {
	ormCategory := domainCategoryToORM(category)
	if err := ormCategory.InsertG(boil.Infer()); err != nil {
		return nil, err
	}
	return ormCategoryToDomain(ormCategory), nil
}

func (repo *PSQLImgRepository) UpdateCategory(category *domain.Category) (*domain.Category, error) {
	ormCategory := domainCategoryToORM(category)
	rows, err := ormCategory.UpdateG(boil.Infer())
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, codes.ErrImgNotFound
	}

	updated, err := repo.FindCategoryByID(ormCategory.ID)
	if err != nil {
		return nil, err
	}

	return updated, err
}

func (repo *PSQLImgRepository) DeleteCategory(id int64) error {
	ormCategory := orm.ImgCategory{
		ID: id,
	}
	rows, err := ormCategory.DeleteG()
	if err != nil {
		return err
	}
	if rows == 0 {
		return codes.ErrImgCategoryNotFound
	}
	return nil
}

func (repo *PSQLImgRepository) ListCategories() ([]*domain.Category, error) {
	ormCategories, err := orm.ImgCategories().AllG()
	if err != nil {
		return nil, err
	}
	return ormCategoriesToDomain(ormCategories), nil
}

func (repo *PSQLImgRepository) FindCategoryByID(id int64) (*domain.Category, error) {
	ormCategory, err := orm.ImgCategories(qm.Where("id = ?", id)).OneG()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrImgCategoryNotFound
		}
		return nil, err
	}
	return ormCategoryToDomain(ormCategory), nil
}

func (repo *PSQLImgRepository) FindCategoryByTitle(title string) (*domain.Category, error) {
	ormCategory, err := orm.ImgCategories(qm.Where(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.Title), title)).OneG()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, codes.ErrImgCategoryNotFound
		}
		return nil, err
	}
	return ormCategoryToDomain(ormCategory), nil
}

func (repo *PSQLImgRepository) CategoryExistByID(id int64) (bool, error) {
	return orm.ImgCategoryExistsG(id)
}

func (repo *PSQLImgRepository) CategoryExistByTitle(title string) (bool, error) {
	exist, err := orm.ImgCategories(qm.Where(fmt.Sprintf("%s = ?", orm.ImgCategoryColumns.Title), title)).ExistsG()
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (repo *PSQLImgRepository) CountCategory() (int64, error) {
	count, err := orm.ImgCategories().CountG()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *PSQLImgRepository) IsCategoryExistImg(id int64) (bool, error) {
	existing2, err := orm.Imgs(qm.Where(fmt.Sprintf("%s = ?", orm.ImgColumns.CategoryID), id), qm.WithDeleted()).ExistsG()
	if err != nil {
		return false, err
	}

	return existing2, nil
}
