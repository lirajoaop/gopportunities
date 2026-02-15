package service

import (
	"fmt"

	"github.com/lirajoaop/gopportunities/repository"
	"github.com/lirajoaop/gopportunities/schemas"
)

type OpeningService struct {
	repo *repository.OpeningRepository
}

func NewOpeningService(repo *repository.OpeningRepository) *OpeningService {
	return &OpeningService{repo: repo}
}

func (s *OpeningService) CreateOpening(opening *schemas.Opening) error {
	return s.repo.Create(opening)
}

func (s *OpeningService) GetOpeningByID(id string) (*schemas.Opening, error) {
	opening, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("opening not found")
	}
	return opening, nil
}

func (s *OpeningService) ListOpenings() ([]schemas.Opening, error) {
	return s.repo.FindAll()
}

func (s *OpeningService) UpdateOpening(id string, role, company, location, link string, remote *bool, salary int64) (*schemas.Opening, error) {
	opening, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("opening not found")
	}

	if role != "" {
		opening.Role = role
	}
	if company != "" {
		opening.Company = company
	}
	if location != "" {
		opening.Location = location
	}
	if remote != nil {
		opening.Remote = *remote
	}
	if link != "" {
		opening.Link = link
	}
	if salary > 0 {
		opening.Salary = salary
	}

	if err := s.repo.Save(opening); err != nil {
		return nil, fmt.Errorf("error updating opening")
	}

	return opening, nil
}

func (s *OpeningService) DeleteOpening(id string) (*schemas.Opening, error) {
	opening, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("opening not found")
	}

	if err := s.repo.Delete(opening); err != nil {
		return nil, fmt.Errorf("error deleting opening")
	}

	return opening, nil
}
