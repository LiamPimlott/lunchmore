package organizations

import (
	"log"

	"github.com/google/uuid"

	"github.com/LiamPimlott/lunchmore/lib/errs"
	"github.com/LiamPimlott/lunchmore/mail"
)

// Service interface to users service
type Service interface {
	Create(o Organization) (Organization, error)
	GetByID(id uint) (Organization, error)
	Invite(i Invitation, inviterID uint) (Invitation, error)
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

	// TODO: retry mechanism
	s.mail.SendInvite(org.Name, i.Code, []string{i.Email})
	if err != nil {
		log.Printf("error sending invitation: %s\n", err)
		return Invitation{}, err
	}

	return inv, nil
}
