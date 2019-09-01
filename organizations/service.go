package organizations

import (
	"log"
)

// Service interface to users service
type Service interface {
	Create(o Organization) (Organization, error)
}

type organizationsService struct {
	repo Repository
}

// NewOrganizationsService will return a struct that implements the organizationsService interface
func NewOrganizationsService(repo Repository) *organizationsService {
	return &organizationsService{
		repo: repo,
	}
}

// Create creates a new organization
func (s *organizationsService) Create(o Organization) (Organization, error) {
	org, err := s.repo.Create(o)
	if err != nil {
		log.Printf("error creating organization: %s\n", err)
		return Organization{}, err
	}

	_, err = s.repo.AddUser(org.ID, o.AdminID)
	if err != nil {
		log.Printf("error creating organization: %s\n", err)
		return Organization{}, err
	}

	return org, nil
}
