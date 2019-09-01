package scheduling

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
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

//GetSchedules returns all schedules
func (r *mysqlSchedulingRepository) GetSchedules() (scheds []Schedule, err error) {

	sql, args, err := sq.Select("*").From("schedules").ToSql()
	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return []Schedule{}, err
	}

	rows, err := r.DB.Query(sql, args)
	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return []Schedule{}, err
	}
	defer rows.Close()

	for rows.Next() {
		s := Schedule{}
		if err := rows.Scan(&s); err != nil {
			log.Printf("error in schedule repo: %s", err.Error())
		}
		scheds = append(scheds, s)
	}

	return scheds, err
}
