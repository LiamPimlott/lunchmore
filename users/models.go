package users

import (
	"github.com/asaskevich/govalidator"
	"github.com/gobuffalo/nulls"
)

// User models an app user
type User struct {
	ID        uint        `json:"id,omitempty" db:"id"`
	OrgID     uint        `json:"org_id,omitempty" db:"org_id"`
	FirstName string      `json:"first_name,omitempty" db:"first_name" valid:"required~first_name is required"`
	LastName  string      `json:"last_name,omitempty" db:"last_name" valid:"required~last_name is required"`
	Email     string      `json:"email,omitempty" db:"email" valid:"required~email is required"`
	Password  string      `json:"password,omitempty" db:"password" valid:"required~password is required"`
	Token     string      `json:"token,omitempty" db:"-"`
	CreatedAt *nulls.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *nulls.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Valid validates a User struct.
func (u User) Valid() (bool, error) {
	return govalidator.ValidateStruct(u)
}

// ValidLogin validates a User struct to be used for logging in.
func (u User) ValidLogin() bool {
	if len(u.Email) <= 0 || len(u.Password) <= 0 {
		return false
	}
	return true
}
