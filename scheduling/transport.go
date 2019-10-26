package scheduling

import (
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"github.com/LiamPimlott/lunchmore/lib/errs"
	"github.com/LiamPimlott/lunchmore/lib/utils"
)

// NewCreateScheduleHandler returns an http handler for creating a schedule.
func NewCreateScheduleHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schedReq := &ScheduleRequest{}

		utils.Decode(w, r, schedReq)

		claims, ok := context.Get(r, "claims").(*utils.CustomClaims)
		if !ok {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		ok, err := schedReq.Valid()
		if !ok || err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, err.Error())
			return
		}

		schedule, err := s.CreateSchedule(*schedReq, claims.ID)
		if err != nil {
			if err.Error() == errs.ErrForbidden.Msg {
				utils.RespondError(w, errs.ErrForbidden.Code, errs.ErrForbidden.Msg, err.Error())
				return
			}
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		utils.Respond(w, schedule)
	}
}

// NewGetOrgSchedulesHandler returns an http handler for retrieving an organization's schedules.
func NewGetOrgSchedulesHandler(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		id, err := strconv.Atoi(params["id"])
		if err != nil {
			utils.RespondError(w, errs.ErrInvalid.Code, errs.ErrInvalid.Msg, "")
			return
		}

		claims, ok := context.Get(r, "claims").(*utils.CustomClaims)
		if !ok {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		if claims.OrgID != uint(id) {
			utils.RespondError(w, errs.ErrForbidden.Code, errs.ErrForbidden.Msg, "")
			return
		}

		schedules, err := s.GetOrgSchedules(uint(id))
		if err != nil {
			utils.RespondError(w, errs.ErrInternal.Code, errs.ErrInternal.Msg, "")
			return
		}

		utils.Respond(w, schedules)
	}
}
