package invites

import (
	"encoding/base64"
	"net/http"

	"github.com/gorilla/context"

	"github.com/LiamPimlott/lunchmore/lib/errs"
	"github.com/LiamPimlott/lunchmore/lib/utils"
)

// NewSendInviteHandler returns an http handler for inviting an email to an organization.
func NewSendInviteHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invReq := &Invitation{}

		utils.Decode(w, r, invReq)

		claims, ok := context.Get(r, "claims").(*utils.CustomClaims)
		if !ok {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		invReq.OrganizationID = claims.OrgID

		ok, err := invReq.Valid()
		if !ok || err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, err.Error())
			return
		}

		invRes, err := s.SendInvite(*invReq, claims.ID)
		if err != nil {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		utils.Respond(w, invRes)
	}
}

// NewAcceptInviteHandler returns an http handler for accepting an invite to an organization.
func NewAcceptInviteHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		joinReq := &JoinRequest{}

		utils.Decode(w, r, joinReq)

		ok, err := joinReq.Valid()
		if !ok || err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, err.Error())
			return
		}

		decodedCode, err := base64.StdEncoding.DecodeString(joinReq.Code)
		if err != nil {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, err.Error())
			return
		}
		joinReq.Code = string(decodedCode)

		user, err := s.AcceptInvite(*joinReq)
		if err != nil {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, err.Error())
			return
		}

		utils.Respond(w, user)
	}
}

// NewGetInviteHandler returns an http handler for getting a sent invite.
func NewGetInviteHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryVals := r.URL.Query()

		code := queryVals.Get("code")
		if code == "" {
			msg := "query param code, is required"
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, msg)
			return
		}

		decodedCode, err := base64.StdEncoding.DecodeString(code)
		if err != nil {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, err.Error())
			return
		}

		orgName, err := s.GetOrgNameByCode(string(decodedCode))
		if err != nil {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, err.Error())
			return
		}

		utils.Respond(w, InviteInfo{OrgName: orgName})
	}
}
