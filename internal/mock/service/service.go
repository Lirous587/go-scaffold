package service

import (
	"scaffold/internal/mock/domain"
)

type service struct {
	repo     domain.MockRepository
}

func NewMockService(repo domain.MockRepository) domain.MockService {
	return &service{
		repo:     repo,
	}
}

func (s *service) Create(mock *domain.Mock) (*domain.Mock, error) {
	return s.repo.Create(mock)
}

func (s *service) Read(id int64) (*domain.Mock, error) {
   return s.repo.FindByID(id)
}

func (s *service) Update(mock *domain.Mock) (*domain.Mock, error) {
	if _, err := s.repo.FindByID(mock.ID); err != nil {
		return nil, err
	}
	return s.repo.Update(mock)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) List(query *domain.MockQuery) (*domain.MockList, error) {
	return s.repo.List(query)
}
