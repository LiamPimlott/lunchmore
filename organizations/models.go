package organizations

import (
	"github.com/asaskevich/govalidator"
	"github.com/gobuffalo/nulls"
)

// Organization models an organization using the app
type Organization struct {
	ID        uint        `json:"id,omitempty" db:"id"`
	AdminID   uint        `json:"admin_id,omitempty" db:"admin_id" valid:"required~admin_id is required"`
	Name      string      `json:"name,omitempty" db:"name" valid:"required~name is required"`
	CreatedAt *nulls.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *nulls.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Valid validates an Org struct.
func (o Organization) Valid() (bool, error) {
	return govalidator.ValidateStruct(o)
}

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
