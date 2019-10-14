package invites

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"log"
)

// Repository interface specifies database api
type Repository interface {
	CreateInvitation(i Invitation) (Invitation, error)
}

type mysqlInvitationsRepository struct {
	DB *sql.DB
}

// NewMysqlOrganizationsRepository returns a struct that implements the Invitations Repository interface
func NewMysqlInvitationsRepository(db *sql.DB) *mysqlInvitationsRepository {
	return &mysqlInvitationsRepository{
		DB: db,
	}
}

// CreateInvitation inserts a new invitation into the db
func (r *mysqlInvitationsRepository) CreateInvitation(i Invitation) (Invitation, error) {
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
