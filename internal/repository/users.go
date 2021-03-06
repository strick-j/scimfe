package repository

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/strick-j/scimfe/internal/model/user"
	"github.com/strick-j/scimfe/internal/service"
	"github.com/strick-j/scimfe/internal/web"
)

const (
	colID       = "id"
	colEmail    = "email"
	colName     = "name"
	colPassword = "password"

	tableUsers = "users"
)

var userCols = []string{colID, colEmail, colName, colPassword}

type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository is UserRepository constructor
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) AllUsers(ctx context.Context) (user.Users, error) {
	q, args, err := psql.Select(userCols...).From(tableUsers).ToSql()
	if err != nil {
		return nil, err
	}

	var out user.Users
	err = r.db.SelectContext(ctx, &out, q, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return out, err
}

func (r UserRepository) AddUser(ctx context.Context, u user.User) (*user.ID, error) {
	q, args, err := psql.Insert(tableUsers).SetMap(map[string]interface{}{
		colEmail:    u.Email,
		colName:     u.Name,
		colPassword: u.PasswordHash,
	}).Suffix("RETURNING " + colID).ToSql()
	if err != nil {
		return nil, err
	}

	newID := new(user.ID)
	return newID, r.db.GetContext(ctx, newID, q, args...)
}

func (r UserRepository) UserByEmail(ctx context.Context, email string) (*user.User, error) {
	q, args, err := psql.Select(userCols...).From(tableUsers).Where(squirrel.Eq{
		colEmail: email,
	}).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	u := new(user.User)
	err = r.db.GetContext(ctx, u, q, args...)
	// UserByEmail is handler differently
	return u, wrapRecordError(err)
}

func (r UserRepository) UserByID(ctx context.Context, uid user.ID) (*user.User, error) {
	q, args, err := psql.Select(userCols...).From(tableUsers).Where(squirrel.Eq{
		colID: uid,
	}).Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	u := new(user.User)
	err = r.db.GetContext(ctx, u, q, args...)
	if err == sql.ErrNoRows {
		return nil, web.NewErrNotFound("user not found")
	}
	return u, err
}

func (r UserRepository) Exists(email string) (bool, error) {
	q, args, err := psql.Select("COUNT(*)").From(tableUsers).Where(squirrel.Eq{
		colEmail: email,
	}).ToSql()
	if err != nil {
		return false, err
	}
	var count uint
	err = r.db.Get(&count, q, args...)
	return count > 0, err
}

func wrapRecordError(err error) error {
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return service.ErrNotExists
	default:
		return err
	}
}
