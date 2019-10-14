package invites

import (
	"github.com/asaskevich/govalidator"
	"github.com/gobuffalo/nulls"
)

// Invitation models an invitation for an email address to join an organization
type Invitation struct {
	ID             uint        `json:"id,omitempty" db:"id"`
	Code           string      `json:"code,omitempty" db:"code"`
	OrganizationID uint        `json:"organization_id,omitempty" db:"organization_id" valid:"required~organization_id is required"`
	Email          string      `json:"email,omitempty" db:"email" valid:"required~email is required"`
	CreatedAt      *nulls.Time `json:"created_at,omitempty" db:"created_at"`
}

// Valid validates an Invitation struct.
func (i Invitation) Valid() (bool, error) {
	return govalidator.ValidateStruct(i)
}

// JoinRequest models a request to accept an organization invite
type JoinRequest struct {
	Code      string `json:"code,omitempty" valid:"required~code is required"`
	FirstName string `json:"first_name,omitempty" valid:"required~first_name is required"`
	LastName  string `json:"last_name,omitempty" valid:"required~last_name is required"`
	Password  string `json:"password,omitempty" valid:"required~password is required"`
}

// Valid validates an JonRequest struct.
func (j JoinRequest) Valid() (bool, error) {
	return govalidator.ValidateStruct(j)
}
