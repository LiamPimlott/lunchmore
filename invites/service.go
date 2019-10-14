package invites

import (
	"log"

	"github.com/google/uuid"

	"github.com/LiamPimlott/lunchmore/lib/errs"
	"github.com/LiamPimlott/lunchmore/mail"
	"github.com/LiamPimlott/lunchmore/organizations"
	"github.com/LiamPimlott/lunchmore/users"
)

// Service interface to users service
type Service interface {
	SendInvite(i Invitation, inviterID uint) (Invitation, error)
	AcceptInvite(j JoinRequest) error
}

type inviteService struct {
	repo  Repository
	mail  mail.Service
	orgs  organizations.Service
	users users.Service
}

// NewInviteService will return a struct that implements the organizationsService interface
func NewInviteService(repo Repository, mail mail.Service, orgs organizations.Service, users users.Service) *inviteService {
	return &inviteService{
		repo:  repo,
		mail:  mail,
		orgs:  orgs,
		users: users,
	}
}

// SendInvite sends an email with an invite link to an organization
func (s *inviteService) SendInvite(i Invitation, inviterID uint) (Invitation, error) {
	org, err := s.orgs.GetByID(i.OrganizationID)
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

func (s *inviteService) AcceptInvite(j JoinRequest) error {
	// get invite

	// create user with provided creds

	// delete invite

	return nil
}
