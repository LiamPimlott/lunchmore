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
	AcceptInvite(j JoinRequest) (users.User, error)
	GetOrgNameByCode(code string) (string, error)
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

func (s *inviteService) AcceptInvite(j JoinRequest) (users.User, error) {
	invite, err := s.repo.GetByCode(j.Code)
	if err != nil {
		log.Printf("error getting invitation by code: %s\n", err)
		return users.User{}, err
	}

	user := users.User{
		OrgID:     invite.OrganizationID,
		FirstName: j.FirstName,
		LastName:  j.LastName,
		Email:     invite.Email,
		Password:  j.Password,
	}

	user, err = s.users.AcceptInvite(user)
	if err != nil {
		log.Printf("error accepting invitation: %s\n", err)
		return users.User{}, err
	}

	err = s.repo.DeleteByID(invite.ID)
	if err != nil {
		log.Printf("error deleting invitation: %s\n", err)
		return users.User{}, err
	}

	return user, nil
}

// GetOrgNameByCode gets the associated organization name by an invite code
func (s *inviteService) GetOrgNameByCode(code string) (string, error) {
	invite, err := s.repo.GetByCode(code)
	if err != nil {
		log.Printf("error getting invitation by code: %s\n", err)
		return "", err
	}

	org, err := s.orgs.GetByID(invite.OrganizationID)
	if err != nil {
		log.Printf("error getting organization by id: %s\n", err)
		return "", err
	}

	return org.Name, nil
}
