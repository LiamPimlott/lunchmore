package organizations

import (
	"log"

	"github.com/LiamPimlott/lunchmore/mail"
)

// Service interface to users service
type Service interface {
	Create(o Organization) (Organization, error)
	GetByID(id uint) (Organization, error)
}

type organizationsService struct {
	repo Repository
	mail mail.Service
}

// NewOrganizationsService will return a struct that implements the organizationsService interface
func NewOrganizationsService(repo Repository, mail mail.Service) *organizationsService {
	return &organizationsService{
		repo: repo,
		mail: mail,
	}
}

// Create creates a new organization
func (s *organizationsService) Create(o Organization) (Organization, error) {
	org, err := s.repo.Create(o)
	if err != nil {
		log.Printf("error creating organization: %s\n", err)
		return Organization{}, err
	}

	return org, nil
}

// GetByID retrieves a organization by its id
func (s *organizationsService) GetByID(id uint) (Organization, error) {
	org, err := s.repo.GetByID(id)
	if err != nil {
		log.Printf("error getting organization: %s\n", err)
		return Organization{}, err
	}

	return org, nil
}
