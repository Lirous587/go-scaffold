package service

import (
	"blog-v4/internal/moment/domain"
)

type service struct {
	repo domain.MomentRepository
}

func NewMomentService(repo domain.MomentRepository) domain.MomentService {
	return &service{
		repo: repo,
	}
}

func (s *service) Get(id int64) (*domain.Moment, error) {
	return s.repo.FindByID(moment)
}

func (s *service) Create(moment *domain.Moment) (*domain.Moment, error) {
	return s.repo.Create(moment)
}

func (s *service) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *service) Update(moment *domain.Moment) (*domain.Moment, error) {
	if _, err := s.repo.FindByID(moment.ID); err != nil {
		return nil, err
	}

	return s.repo.Update(moment)
}

func (s *service) ListMoments(query *domain.MomentQuery) (*domain.MomentPages, error) {
	return s.repo.ListMoments(query)
}

func (s *service) RandomN(count int8) ([]*domain.Moment, error) {
	return s.repo.RandomN(count)
}
