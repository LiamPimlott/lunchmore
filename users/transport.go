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

// NewCreateUserHandler returns an http handler for signing up a new user and organization.
func NewCreateUserHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userReq := &User{}

		utils.Decode(w, r, userReq)

		ok, err := userReq.Valid()
		if !ok || err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, err.Error())
			return
		}

		userRes, err := s.Create(*userReq)
		if err != nil {
			log.Println("error creating user")
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		utils.Respond(w, userRes)
	}
}

// NewLoginHandler returns an http handler for logging in users
func NewLoginHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := &User{}

		utils.Decode(w, r, body)

		ok := body.ValidLogin()
		if !ok {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, "")
			return
		}

		tkn, err := s.Login(*body)
		if err != nil {
			log.Printf("error logging in user: %s\n", err)
			if err == sql.ErrNoRows {
				utils.RespondError(w, errs.ErrNotFound.Code, errs.ErrNotFound.Msg, "")
				return
			}
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
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
