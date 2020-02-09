package users

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/LiamPimlott/lunchmore/lib/utils"
	orgs "github.com/LiamPimlott/lunchmore/organizations"
)

// Service interface to users service
type Service interface {
	Signup(sr SignupRequest) (User, error)
	Login(u User) (User, error)
	AcceptInvite(u User) (User, error)
	GetUsersMap(usrIDs []uint) (map[uint]User, error)
	GetByID(idRequested, idClaimed uint) (User, error)
	RefreshJwt(usrID uint) (string, error)
}

type usersService struct {
	repo   Repository
	o      orgs.Service
	secret string
}

// NewUsersService will return a struct that implements the UsersService interface
func NewUsersService(repo Repository, o orgs.Service, secret string) *usersService {
	return &usersService{
		repo:   repo,
		o:      o,
		secret: secret,
	}
}

//Signup creates a new user with a new organization and issues a token
func (s *usersService) Signup(sr SignupRequest) (User, error) {
	// TODO santitize and validate input
	pass, err := bcrypt.GenerateFromPassword([]byte(sr.Password), bcrypt.MinCost)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return User{}, err
	}
	sr.Password = string(pass)

	usr, err := s.repo.Create(User{
		FirstName: sr.FirstName,
		LastName:  sr.LastName,
		Email:     sr.Email,
		Password:  sr.Password,
	})
	if err != nil {
		log.Printf("error creating user: %s\n", err)
		return User{}, err
	}
	log.Printf("SR: %+v\n", sr)

	org, err := s.o.Create(orgs.Organization{AdminID: usr.ID, Name: sr.OrgName})
	if err != nil {
		log.Printf("error creating user: %s\n", err)
		return User{}, err
	}

	err = s.repo.UpdateOrganization(usr.ID, org.ID)
	if err != nil {
		log.Printf("error creating user: %s\n", err)
		return User{}, err
	}
	usr.OrgID = org.ID

	token, err := utils.GenerateToken(usr.ID, org.ID, s.secret)
	if err != nil {
		log.Printf("error creating user: %s\n", err)
		return User{}, err
	}
	usr.Token = token

	return usr, nil
}

// Login validates an email/pass and return a token
func (s *usersService) Login(u User) (User, error) {
	usr, err := s.repo.GetByEmail(u.Email)
	if err != nil {
		log.Printf("error logging in user: %s\n", err)
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(u.Password))
	if err != nil {
		log.Printf("error logging in user: %s\n", err)
		return User{}, err
	}

	token, err := utils.GenerateToken(usr.ID, usr.OrgID, s.secret)
	if err != nil {
		log.Printf("error logging in user: %s\n", err)
		return User{}, err
	}

	usr.Password = ""
	usr.Token = token

	return usr, nil
}

// GetUsers retrieves a map of users by id
func (s *usersService) GetUsersMap(usrIDs []uint) (map[uint]User, error) {
	usrs, err := s.repo.GetUsers(usrIDs)
	if err != nil {
		log.Printf("error getting users: %s\n", err)
		return map[uint]User{}, err
	}

	usrsMap := map[uint]User{}
	for _, u := range usrs {
		usrsMap[u.ID] = u
	}

	return usrsMap, nil
}

// GetByID retrieves a user by their id
func (s *usersService) GetByID(idRequested, idClaimed uint) (User, error) {

	usr, err := s.repo.GetByID(idRequested)
	if err != nil {
		log.Printf("error getting user: %s\n", err)
		return User{}, err
	}

	if idRequested != idClaimed {
		usr = User{
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
		}
	}

	return usr, nil
}

//Signup creates a new user with a new organization and issues a token
func (s *usersService) AcceptInvite(u User) (User, error) {
	// TODO santitize and validate input
	passEncrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		log.Printf("error generating password hash: %s", err.Error())
		return User{}, err
	}
	u.Password = string(passEncrypted)

	u, err = s.repo.Create(u)
	if err != nil {
		log.Printf("error calling user repo: %s\n", err)
		return User{}, err
	}

	token, err := utils.GenerateToken(u.ID, u.OrgID, s.secret)
	if err != nil {
		log.Printf("error generating token: %s\n", err)
		return User{}, err
	}
	u.Token = token

	return u, nil
}

// RefreshJwt creates a fresh jwt for a user
func (s *usersService) RefreshJwt(usrID uint) (string, error) {
	usr, err := s.repo.GetByID(usrID)
	if err != nil {
		log.Printf("error getting user by id: %s\n", err)
		return "", err
	}

	tkn, err := utils.GenerateToken(usr.ID, usr.OrgID, s.secret)
	if err != nil {
		log.Printf("error generating token: %s\n", err)
		return "", err
	}

	return tkn, nil
}
