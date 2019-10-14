package scheduling

import (
	"github.com/asaskevich/govalidator"
	"github.com/gobuffalo/nulls"
)

// Schedule models a cron configuration for scheduling the matching process
type Schedule struct {
	ID        uint        `json:"id,omitempty" db:"id"`
	OrgID     uint        `json:"org_id,omitempty" db:"org_id"`
	Spec      string      `json:"config_string,omitempty" db:"config_string"`
	CreatedAt *nulls.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *nulls.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Valid validates a Schedule struct.
func (s Schedule) Valid() (bool, error) {
	return govalidator.ValidateStruct(s)
}

// ScheduleUser represents a user's inclusion in a schedule
type ScheduleUser struct {
	ID         uint        `json:"id,omitempty" db:"id"`
	UserID     uint        `json:"user_id,omitempty" db:"user_id"`
	ScheduleID uint        `json:"schedule_id,omitempty" db:"schedule_id"`
	CreatedAt  *nulls.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  *nulls.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Valid validates a ScheduleUser struct.
func (s ScheduleUser) Valid() (bool, error) {
	return govalidator.ValidateStruct(s)
}

type LunchMatch struct {
	ID         uint        `json:"id,omitempty" db:"id"`
	UserID1    uint        `json:"user_id_1,omitempty" db:"user_id_1"`
	UserID2    uint        `json:"user_id_2,omitempty" db:"user_id_2"`
	ScheduleID uint        `json:"schedule_id,omitempty" db:"schedule_id"`
	CreatedAt  *nulls.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  *nulls.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Valid validates a LunchMatch struct.
func (l LunchMatch) Valid() (bool, error) {
	return govalidator.ValidateStruct(l)
}

// type ScheduleDay string

// const (
// 	Sunday    ScheduleDay = "SUN"
// 	Monday    ScheduleDay = "MON"
// 	Tuesday   ScheduleDay = "TUE"
// 	Wednesday ScheduleDay = "WED"
// 	Thursday  ScheduleDay = "THU"
// 	Friday    ScheduleDay = "FRI"
// 	Saturday  ScheduleDay = "SAT"
// )

// ScheduleRequest models a request for a new schedule
type ScheduleRequest struct {
	OrgID uint     `json:"org_id,omitempty"`
	Days  []string `json:"days,omitempty" valid:"required~field days is required, in(SUN|MON|TUE|WED|THU|FRI|SAT)"`
	// Interval string   `json:"interval,omitempty" valid:"required~field interval is required`
}

// Valid validates a ScheduleRequest struct.
func (sr ScheduleRequest) Valid() (bool, error) {
	return govalidator.ValidateStruct(sr)
}
