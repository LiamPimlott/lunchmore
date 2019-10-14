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
