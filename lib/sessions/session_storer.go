package sessions

import (
	"log"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"github.com/LiamPimlott/lunchmore/lib/errs"
)

// SessionStorer interface for interacting with a seesion store
type SessionStorer interface {
	Set(r *http.Request, w http.ResponseWriter, usrID uint) error
	Validate(r *http.Request) (uint, error)
}

type sessionStorer struct {
	store sessions.Store
	name  string
}

// NewSessionStorer constructs a SessionStorer
func NewSessionStorer(store sessions.Store, name string) *sessionStorer {
	return &sessionStorer{
		store,
		name,
	}
}

// Set sets a new or existing session
func (s *sessionStorer) Set(r *http.Request, w http.ResponseWriter, usrID uint) error {
	session, err := s.store.Get(r, s.name)
	if err != nil && err != securecookie.ErrMacInvalid {
		log.Printf("failed to get session: %s", err.Error())
		return errs.ErrInternal
	}

	session.Values["authenticated"] = true
	session.Values["userID"] = usrID
	if err := s.store.Save(r, w, session); err != nil {
		log.Printf("failed to save session: %s", err.Error())
		return errs.ErrInternal
	}

	return nil
}

// Validate checks if a request contains a valid session
func (s *sessionStorer) Validate(r *http.Request) (uint, error) {
	session, err := s.store.Get(r, s.name)
	if err != nil {
		log.Printf("failed to get session: %s", err.Error())
		return 0, errs.ErrForbidden
	}

	if a, ok := session.Values["authenticated"].(bool); session.IsNew || !ok || !a {
		log.Printf("session un-authenticated")
		return 0, errs.ErrForbidden
	}

	usrID, ok := session.Values["userID"].(uint)
	if !ok || usrID == 0 {
		log.Printf("session is missing userID")
		return 0, errs.ErrForbidden
	}

	return usrID, nil
}
