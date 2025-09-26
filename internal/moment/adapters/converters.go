package adapters

import (
	"blog-v4/internal/common/orm"
	"blog-v4/internal/moment/domain"
	"github.com/aarondl/null/v8"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/types/pgeo"
)

func domainMomentToORM(moment *domain.Moment) *orm.Moment {
	if moment == nil {
		return nil
	}

	// 非null项
	ormMoment := &orm.Moment{
		ID:          moment.ID,
		Title:       moment.Title,
		Content:     moment.Content,
		CreatedAt:   moment.CreatedAt,
		UpdatedAt:   moment.UpdatedAt,
		Location:    null.String{Valid: false},
		Coordinates: pgeo.NullPoint{Valid: false},
		CoverURL:    moment.CoverURL,
	}

	// 处理null项
	if moment.Location != "" {
		ormMoment.Location = null.StringFrom(moment.Location)
	}

	if moment.Coordinates != nil && moment.Coordinates.IsValid() {
		point := pgeo.NewPoint(moment.Coordinates.X, moment.Coordinates.Y)
		ormMoment.Coordinates = pgeo.NewNullPoint(point, true)
	}

	return ormMoment
}

func ormMomentToDomain(ormMoment *orm.Moment) *domain.Moment {
	if ormMoment == nil {
		return nil
	}

	moment := &domain.Moment{
		// 非null项
		ID:        ormMoment.ID,
		Title:     ormMoment.Title,
		Content:   ormMoment.Content,
		CreatedAt: ormMoment.CreatedAt,
		UpdatedAt: ormMoment.UpdatedAt,
		CoverURL:  ormMoment.CoverURL,
	}

	// 处理null项
	if ormMoment.Location.Valid {
		moment.Location = ormMoment.Location.String
	}

	if ormMoment.Coordinates.Valid {
		point := ormMoment.Coordinates.Point
		moment.Coordinates = &domain.MomentCoordinate{
			X: point.X,
			Y: point.Y,
		}
	}
	return moment
}

func ormMomentsToDomain(ormMoments []*orm.Moment) []*domain.Moment {
	if len(ormMoments) == 0 {
		return nil
	}

	moments := make([]*domain.Moment, 0, len(ormMoments))
	for _, ormMoment := range ormMoments {
		if ormMoment != nil {
			moments = append(moments, ormMomentToDomain(ormMoment))
		}
	}
	return moments
}
