package users

import (
	"database/sql"
	"log"

	sq "github.com/Masterminds/squirrel"
)

// Repository interface specifies database api
type Repository interface {
	Create(u User) (User, error)
	GetPassword(email string) (User, error)
	GetUsers(usrIDs []uint) ([]User, error)
	GetById(id int) (User, error)
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
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"email":      u.Email,
		"password":   u.Password,
	}).ToSql()

	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	res, err := r.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	return User{ID: uint(id)}, nil
}

// GetPassword retrieves the id, email and password for an email
func (r *mysqlUsersRepository) GetPassword(email string) (User, error) {
	var usr User

	stmnt, args, err := sq.Select("id", "email", "password").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	err = r.DB.QueryRow(stmnt, args...).Scan(&usr.ID, &usr.Email, &usr.Password)
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

// GetById get user by id
func (r *mysqlUsersRepository) GetById(id int) (User, error) {
	var usr User

	stmnt, args, err := sq.Select("id", "first_name", "last_name", "email").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		log.Printf("error in user repo: %s", err.Error())
		return User{}, err
	}

	err = r.DB.QueryRow(stmnt, args...).Scan(
		&usr.ID,
		&usr.FirstName,
		&usr.LastName,
		&usr.Email,
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
