package organizations

import (
	"github.com/LiamPimlott/lunchmore/lib/errs"
	"github.com/google/uuid"
	"log"
)

// Service interface to users service
type Service interface {
	Create(o Organization) (Organization, error)
	GetByID(id uint) (Organization, error)
	Invite(i Invitation, inviterID uint) (Invitation, error)
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

// Invite sends an email with an invite link to an organization
func (s *organizationsService) Invite(i Invitation, inviterID uint) (Invitation, error) {
	org, err := s.GetByID(i.OrganizationID)
	if err != nil {
		log.Printf("error getting organization: %s\n", err)
		return Invitation{}, err
	}

	if org.AdminID != inviterID {
		log.Printf("error admin id does not match identity: %s\n", errs.ErrForbidden)
		return Invitation{}, errs.ErrForbidden
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("error creating invitation code: %s\n", err)
		return Invitation{}, err
	}

	code := uuid.String()
	if code == "" {
		log.Printf("error creating invitation code: %s\n", err)
		return Invitation{}, err
	}
	i.Code = code

	inv, err := s.repo.CreateInvitation(i)
	if err != nil {
		log.Printf("error creating invitation: %s\n", err)
		return Invitation{}, err
	}

	// TODO: send email with link and bas64 invite code

	return inv, nil
}
