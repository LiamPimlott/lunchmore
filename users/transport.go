package users

import (
	"database/sql"
	"log"
	"net/http"
	// "strconv"

	// "github.com/gorilla/context"
	// "github.com/gorilla/mux"

	"github.com/LiamPimlott/lunchmore/lib/errs"
	"github.com/LiamPimlott/lunchmore/lib/utils"
)

// NewSignupHandler returns an http handler for signing up a new user and organization.
func NewSignupHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signupReq := &SignupRequest{}

		utils.Decode(w, r, signupReq)

		ok, err := signupReq.Valid()
		if !ok || err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, err.Error())
			return
		}

		usr, err := s.Signup(*signupReq)
		if err != nil {
			log.Printf("error creating user: %s", err)
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		utils.Respond(w, usr)
	}
}

// NewLoginHandler returns an http handler for logging in users
func NewLoginHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := &User{}

		utils.Decode(w, r, body)

		ok := body.ValidLogin()
		if !ok {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, "Invalid request data.")
			return
		}

		tkn, err := s.Login(*body)
		if err != nil {
			log.Printf("error logging in user: %s\n", err)
			if err == sql.ErrNoRows {
				utils.RespondError(w, errs.ErrNotFound.Code, errs.ErrNotFound.Msg, "Email or password is incorrect.")
				return
			}
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "An error has occured.")
			return
		}

		utils.Respond(w, tkn)
	}
}

// NewGetUserByIDHandler returns an http handler for getting users by id
// func NewGetUserByIDHandler(s Service) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		rtPrms := mux.Vars(r)

// 		idStrng, ok := rtPrms["id"]
// 		if !ok {
// 			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, "missing user id in url")
// 			return
// 		}

// 		idRequested, err := strconv.Atoi(idStrng)
// 		if err != nil {
// 			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, "invalid user id in url")
// 			return
// 		}

// 		claims, ok := context.Get(r, "claims").(*utils.CustomClaims)
// 		if !ok {
// 			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
// 			return
// 		}

// 		usr, err := s.GetByID(idRequested, int(claims.ID))
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				utils.RespondError(w, errs.ErrNotFound.Code, errs.ErrNotFound.Msg, "")
// 				return
// 			}
// 			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
// 			return
// 		}

// 		utils.Respond(w, usr)
// 	}
// }
