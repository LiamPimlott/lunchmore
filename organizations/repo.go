package organizations

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"log"
)

// Repository interface specifies database api
type Repository interface {
	AddUser(oID, uID uint) (uint, error)
	Create(o Organization) (Organization, error)
}

type mysqlOrganizationsRepository struct {
	DB *sql.DB
}

// NewMysqlOrganizationsRepository returns a struct that implements the OrganizationsRepository interface
func NewMysqlOrganizationsRepository(db *sql.DB) *mysqlOrganizationsRepository {
	return &mysqlOrganizationsRepository{
		DB: db,
	}
}

// AddUser adds a org/user pair to the organization_users table
func (r *mysqlOrganizationsRepository) AddUser(oID, uID uint) (uint, error) {
	sql, args, err := sq.Insert("organization_users").SetMap(sq.Eq{
		"organization_id": oID,
		"user_id":         uID,
	}).ToSql()

	if err != nil {
		log.Printf("error in organization repo: %s", err.Error())
		return 0, err
	}

	res, err := r.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("error in organization repo: %s", err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error in organization repo: %s", err.Error())
		return 0, err
	}

	return uint(id), nil
}

// Create inserts a new organization into the db
func (r *mysqlOrganizationsRepository) Create(o Organization) (Organization, error) {
	// TODO: validate & sanitize

	sql, args, err := sq.Insert("organizations").SetMap(sq.Eq{
		"admin_id": o.AdminID,
		"name":     o.Name,
	}).ToSql()

	if err != nil {
		log.Printf("error in organization repo: %s", err.Error())
		return Organization{}, err
	}

	res, err := r.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("error in organization repo: %s", err.Error())
		return Organization{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error in organization repo: %s", err.Error())
		return Organization{}, err
	}

	return Organization{ID: uint(id)}, nil
}
