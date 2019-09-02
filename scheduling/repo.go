package scheduling

import (
	"database/sql"
	"errors"
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

// GetSchedules returns all schedules
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

// GetScheduleUsers gets all users part of a schedule
func (r *mysqlSchedulingRepository) GetScheduleUsers(sID uint) (su []ScheduleUser, err error) {

	sql, args, err := sq.
		Select("*").
		From("schedules_users").
		Where(sq.Eq{"id": sID}).ToSql()

	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return su, err
	}

	rows, err := r.DB.Query(sql, args)
	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return su, err
	}
	defer rows.Close()

	for rows.Next() {
		s := ScheduleUser{}
		if err := rows.Scan(&s); err != nil {
			log.Printf("error in schedule repo: %s", err.Error())
		}
		su = append(su, s)
	}
	return su, err
}

func (r *mysqlSchedulingRepository) SaveLunchMatches(lm []LunchMatch) error {

	stmnt := sq.Insert("lunch_matches").Columns(
		"user_id_1",
		"user_id_1",
		"schedule_id",
	)

	for _, m := range lm {
		stmnt.Values(m.UserID1, m.UserID2, m.ScheduleID)
	}

	sql, args, err := stmnt.ToSql()
	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return err
	}

	res, err := r.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return err
	}

	numRows, err := res.RowsAffected()
	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return err
	}
	if int(numRows) != len(lm) {
		err := errors.New("error: rows affected does not match len(m)")
		log.Printf("%s", err.Error())
		return err
	}

	return nil
}
