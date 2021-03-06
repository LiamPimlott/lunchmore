package scheduling

import (
	"database/sql"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
)

// Repository interface specifies database api
type Repository interface {
	CreateSchedule(s Schedule) (Schedule, error)
	GetOrgSchedules(ordID uint) ([]Schedule, error)
	GetByID(id uint) (Schedule, error)
	GetSchedules() ([]Schedule, error)
	CreateScheduleUser(s ScheduleUser) (ScheduleUser, error)
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

// CreateSchedule insert schedule into database
func (r *mysqlSchedulingRepository) CreateSchedule(s Schedule) (Schedule, error) {
	sql, args, err := sq.Insert("schedules").SetMap(sq.Eq{
		"org_id": s.OrgID,
		"spec":   s.Spec,
	}).ToSql()

	if err != nil {
		log.Printf("error assembling create schedule statement: %s", err.Error())
		return Schedule{}, err
	}

	res, err := r.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("error executing create schedule statement: %s", err.Error())
		return Schedule{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error accesing last inserted schedule id: %s", err.Error())
		return Schedule{}, err
	}
	s.ID = uint(id)

	return s, nil
}

// GetByID get schedule by id
func (r *mysqlSchedulingRepository) GetByID(id uint) (Schedule, error) {
	var sched Schedule

	stmnt, args, err := sq.
		Select("id", "org_id", "spec").
		From("schedules").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		log.Printf("error in build get schedule by id statement: %s", err.Error())
		return Schedule{}, err
	}

	err = r.DB.QueryRow(stmnt, args...).Scan(
		&sched.ID,
		&sched.OrgID,
		&sched.Spec,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("schedule %d not found.", sched.ID)
			return Schedule{}, err
		}
		log.Printf("error executing get schedule by id query: %s", err.Error())
		return Schedule{}, err
	}

	return sched, nil
}

// GetOrgSchedules gets schedules by organization
func (r *mysqlSchedulingRepository) GetOrgSchedules(orgID uint) ([]Schedule, error) {
	sql, args, err := sq.
		Select("*").
		From("schedules").
		Where(sq.Eq{"org_id": orgID}).
		ToSql()
	if err != nil {
		log.Printf("error creating get org schedules query: %s", err.Error())
		return []Schedule{}, err
	}

	rows, err := r.DB.Query(sql, args...)
	if err != nil {
		log.Printf("error executing get org schedules query: %s", err.Error())
		return []Schedule{}, err
	}
	defer rows.Close()

	scheds, err := scanSchedules(rows)
	if err != nil {
		log.Printf("error scanning schedules: %s", err.Error())
		return []Schedule{}, err
	}

	return scheds, err
}

// GetSchedules returns all schedules
func (r *mysqlSchedulingRepository) GetSchedules() (scheds []Schedule, err error) {

	sql, args, err := sq.Select("*").From("schedules").ToSql()
	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return []Schedule{}, err
	}

	rows, err := r.DB.Query(sql, args...)
	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return []Schedule{}, err
	}
	defer rows.Close()

	scheds, err = scanSchedules(rows)
	if err != nil {
		log.Printf("error scanning schedules: %s", err.Error())
		return []Schedule{}, err
	}

	return scheds, err
}

// CreateSheduleUser creates a schedule user
func (r *mysqlSchedulingRepository) CreateScheduleUser(s ScheduleUser) (ScheduleUser, error) {
	sql, args, err := sq.Insert("schedule_users").SetMap(sq.Eq{
		"user_id":     s.UserID,
		"schedule_id": s.ScheduleID,
	}).ToSql()

	if err != nil {
		log.Printf("error assembling create schedule statement: %s", err.Error())
		return ScheduleUser{}, err
	}

	res, err := r.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("error executing create schedule statement: %s", err.Error())
		return ScheduleUser{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error accesing last inserted schedule id: %s", err.Error())
		return ScheduleUser{}, err
	}
	s.ID = uint(id)

	return s, nil
}

// GetScheduleUsers gets all users part of a schedule
func (r *mysqlSchedulingRepository) GetScheduleUsers(sID uint) (su []ScheduleUser, err error) {

	sql, args, err := sq.
		Select("*").
		From("schedule_users").
		Where(sq.Eq{"schedule_id": sID}).ToSql()

	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return su, err
	}

	rows, err := r.DB.Query(sql, args...)
	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return su, err
	}
	defer rows.Close()

	for rows.Next() {
		s := ScheduleUser{}
		dest := []interface{}{
			&s.ID,
			&s.UserID,
			&s.ScheduleID,
			&s.CreatedAt,
			&s.UpdatedAt,
		}
		if err := rows.Scan(dest...); err != nil {
			log.Printf("error in schedule repo: %s", err.Error())
		}
		su = append(su, s)
	}

	return su, err
}

func (r *mysqlSchedulingRepository) SaveLunchMatches(lm []LunchMatch) error {

	stmnt := sq.Insert("lunch_matches").Columns(
		"user_id_1",
		"user_id_2",
		"schedule_id",
	)

	for _, m := range lm {
		stmnt = stmnt.Values(m.UserID1, m.UserID2, m.ScheduleID)
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

func scanSchedules(rows *sql.Rows) ([]Schedule, error) {
	var scheds []Schedule

	for rows.Next() {
		s := Schedule{}
		dest := []interface{}{
			&s.ID,
			&s.OrgID,
			&s.Spec,
			&s.CreatedAt,
			&s.UpdatedAt,
		}
		if err := rows.Scan(dest...); err != nil {
			log.Printf("error scanning schedules in schedule repo: %s", err.Error())
			return scheds, err
		}
		scheds = append(scheds, s)
	}

	return scheds, nil
}
