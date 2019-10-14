package organizations

import (
	"log"
	"net/http"

	"github.com/gorilla/context"

	"github.com/LiamPimlott/lunchmore/lib/errs"
	"github.com/LiamPimlott/lunchmore/lib/utils"
)

// NewCreateOrganizationHandler returns an http handler for signing up a organization.
func NewCreateOrganizationHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgReq := &Organization{}

		utils.Decode(w, r, orgReq)

		claims, ok := context.Get(r, "claims").(*utils.CustomClaims)
		if !ok {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		orgReq.AdminID = claims.ID

		ok, err := orgReq.Valid()
		if !ok || err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, err.Error())
			return
		}

		orgRes, err := s.Create(*orgReq)
		if err != nil {
			log.Println("error creating organization")
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		utils.Respond(w, orgRes)
	}
}
