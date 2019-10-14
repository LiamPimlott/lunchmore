package invites

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
)

// Repository interface specifies database api
type Repository interface {
	CreateInvitation(i Invitation) (Invitation, error)
	GetByCode(code string) (Invitation, error)
	DeleteByID(id uint) error
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

// GetByCode attempts to retrieve an invitation by uuid string
func (r *mysqlInvitationsRepository) GetByCode(code string) (Invitation, error) {
	var inv Invitation

	stmnt, args, err := sq.
		Select("*").
		From("invitations").
		Where(sq.Eq{"code": code}).
		ToSql()
	if err != nil {
		log.Printf("error in organizations repo: %s", err.Error())
		return Invitation{}, err
	}

	err = r.DB.QueryRow(stmnt, args...).Scan(
		&inv.ID,
		&inv.Code,
		&inv.OrganizationID,
		&inv.Email,
		&inv.CreatedAt,
		&inv.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("invitation %s not found.", code)
			return Invitation{}, err
		}
		log.Printf("error in invitations repo: %s", err.Error())
		return Invitation{}, err
	}

	return inv, nil
}

// DeleteByID attempts to delete an invitation by id
func (r *mysqlInvitationsRepository) DeleteByID(id uint) error {
	stmnt, args, err := sq.
		Delete("invitations").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		log.Printf("error preparing delete statement: %s", err.Error())
		return err
	}

	res, err := r.DB.Exec(stmnt, args...)
	if err != nil {
		log.Printf("error executing statement: %s", err.Error())
		return err
	}

	numAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("error checking rows affected: %s", err.Error())
		return err
	} else if numAffected != 1 {
		err := fmt.Errorf("incorrect number of rows affected: %d", numAffected)
		log.Printf("error checking rows affected: %s", err.Error())
		return err
	}

	return nil
}
