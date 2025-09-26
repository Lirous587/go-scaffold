package adapters

import (
	"scaffold/internal/common/orm"
	"scaffold/internal/mock/domain"
	"github.com/aarondl/null/v8"
)

func domainMockToORM(mock *domain.Mock) *orm.Mock {
	if mock == nil {
		return nil
	}

	// 非null项
	ormMock := &orm.Mock{
		ID:        		mock.ID,
        CreatedAt: 		mock.CreatedAt,
        UpdatedAt: 		mock.UpdatedAt,
	}

	// 处理null项
	if mock.Description != "" {
	 	ormMock.Description = null.StringFrom(mock.Description)
	}

	return ormMock
}

func ormMockToDomain(ormMock *orm.Mock) *domain.Mock {
	if ormMock == nil {
		return nil
	}

	// 非null项
	mock := &domain.Mock{
		ID:        		ormMock.ID,
		CreatedAt: 		ormMock.CreatedAt,
		UpdatedAt: 		ormMock.UpdatedAt,
	}

	// 处理null项
	if ormMock.Description.Valid {
 	 	mock.Description = ormMock.Description.String
	}

	return mock
}

func ormMocksToDomain(ormMocks []*orm.Mock) []*domain.Mock {
	if len(ormMocks) == 0 {
		return nil
	}

	mocks := make([]*domain.Mock, 0, len(ormMocks))
	for _, ormMock := range ormMocks {
		if ormMock != nil {
			mocks = append(mocks, ormMockToDomain(ormMock))
		}
	}
	return mocks
}

