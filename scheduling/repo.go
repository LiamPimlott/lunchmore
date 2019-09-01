package scheduling

import (
	"database/sql"
	// sq "github.com/Masterminds/squirrel"
	// "log"
)

// Repository interface specifies database api
type Repository interface {
	GetSchedules() ([]Schedule, error)
	GetScheduleUsers(sID uint) ([]ScheduleUser, error)
	SaveLunchMatches(lm []LunchMatch) error
}

type mysqlSchedulingRepository struct {
	DB *sql.DB
}

// NewMysqlSchedulingRepository returns a struct that implements the SchedulingRepository interface
func NewMysqlSchedulingRepository(db *sql.DB) *mysqlSchedulingRepository {
	return &mysqlSchedulingRepository{
		DB: db,
	}
}
