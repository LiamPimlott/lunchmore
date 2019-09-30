package organizations

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"log"
)

// Repository interface specifies database api
type Repository interface {
	Create(o Organization) (Organization, error)
	GetByID(id uint) (Organization, error)
	CreateInvitation(i Invitation) (Invitation, error)
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

// GetByID get organization by id
func (r *mysqlOrganizationsRepository) GetByID(id uint) (Organization, error) {
	var org Organization

	stmnt, args, err := sq.
		Select("id", "admin_id", "name").
		From("organizations").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		log.Printf("error in organizations repo: %s", err.Error())
		return Organization{}, err
	}

	err = r.DB.QueryRow(stmnt, args...).Scan(
		&org.ID,
		&org.AdminID,
		&org.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("organization %d not found.", org.ID)
			return Organization{}, err
		}
		log.Printf("error in organizations repo: %s", err.Error())
		return Organization{}, err
	}

	return org, nil
}

// CreateInvitation inserts a new invitation into the db
func (r *mysqlOrganizationsRepository) CreateInvitation(i Invitation) (Invitation, error) {
	// TODO: validate & sanitize

	sql, args, err := sq.Insert("invitations").SetMap(sq.Eq{
		"organization_id": i.OrganizationID,
		"email":           i.Email,
		"code":            i.Code,
	}).ToSql()

	if err != nil {
		log.Printf("error in organization repo: %s", err.Error())
		return Invitation{}, err
	}

	res, err := r.DB.Exec(sql, args...)
	if err != nil {
		log.Printf("error in organization repo: %s", err.Error())
		return Invitation{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("error in organization repo: %s", err.Error())
		return Invitation{}, err
	}

	return Invitation{ID: uint(id)}, nil
}
