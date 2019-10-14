package users

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
)

// Repository interface specifies database api
type Repository interface {
	Create(u User) (User, error)
	GetByEmail(email string) (User, error)
	GetUsers(usrIDs []uint) ([]User, error)
	GetByID(id uint) (User, error)
	UpdateOrganization(uID, oID uint) error
}

type mysqlUsersRepository struct {
	DB *sql.DB
}

// NewMysqlUsersRepository returns a struct that implements the UsersRepository interface
func NewMysqlUsersRepository(db *sql.DB) *mysqlUsersRepository {
	return &mysqlUsersRepository{
		DB: db,
	}
}

// Create inserts a new user into the db
func (r *mysqlUsersRepository) Create(u User) (User, error) {
	// TODO: validate & sanitize

	sql, args, err := sq.Insert("users").SetMap(sq.Eq{
		"organization_id": u.OrgID,
		"first_name":      u.FirstName,
		"last_name":       u.LastName,
		"email":           u.Email,
		"password":        u.Password,
	}).ToSql()

	if err != nil {
		log.Printf("error assembling create user statement: %s", err.Error())
		return User{}, err
	}

	res, err := r.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("error executing create user statement: %s", err.Error())
		return User{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error accesing last inserted user id: %s", err.Error())
		return User{}, err
	}

	usr, err := r.GetByID(uint(id))
	if err != nil {
		log.Printf("error retrieving created user by id: %s", err.Error())
		return User{}, err
	}

	return usr, nil
}

// GetByEmail retrieves a user by email
func (r *mysqlUsersRepository) GetByEmail(email string) (User, error) {
	var usr User

	stmnt, args, err := sq.Select(
		"id", "organization_id", "first_name",
		"last_name", "email", "password",
	).
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	err = r.DB.QueryRow(stmnt, args...).Scan(
		&usr.ID,
		&usr.OrgID,
		&usr.FirstName,
		&usr.LastName,
		&usr.Email,
		&usr.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("user %d not found.", usr.ID)
			return User{}, err
		}
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	return usr, nil
}

// GetUsers get users by a list of ids
func (r *mysqlUsersRepository) GetUsers(usrIDs []uint) (usrs []User, err error) {

	sql, args, err := sq.
		Select("id", "first_name", "last_name", "email").
		From("users").
		Where(sq.Eq{"id": usrIDs}).ToSql()

	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return usrs, err
	}

	rows, err := r.DB.Query(sql, args...)
	if err != nil {
		log.Printf("error in schedule repo: %s", err.Error())
		return usrs, err
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}

		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			log.Printf("error in schedule repo: %s", err.Error())
		}

		usrs = append(usrs, u)
	}

	return usrs, err
}

// GetByID get user by id excluding password
func (r *mysqlUsersRepository) GetByID(id uint) (User, error) {
	var usr User

	stmnt, args, err := sq.Select(
		"id", "organization_id", "first_name",
		"last_name", "email",
	).
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		log.Printf("error in assembling get by id statement: %s", err.Error())
		return User{}, err
	}

	err = r.DB.QueryRow(stmnt, args...).Scan(
		&usr.ID,
		&usr.OrgID,
		&usr.FirstName,
		&usr.LastName,
		&usr.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("user %d not found.", usr.ID)
			return User{}, err
		}
		log.Printf("error executing query: %s", err.Error())
		return User{}, err
	}

	return usr, nil
}

// UpdateOrganization adds an organization id to a user
func (r *mysqlUsersRepository) UpdateOrganization(uID, oID uint) error {
	stmnt, args, err := sq.Update("users").
		Set("organization_id", oID).
		Where(sq.Eq{"id": uID}).
		ToSql()
	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return err
	}

	result, err := r.DB.Exec(stmnt, args...)
	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil || affected != 1 {
		log.Printf("error in user repo: %s", err.Error())
		return err
	}

	return nil
}
