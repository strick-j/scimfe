package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/strick-j/scimfe/internal/model"
	"github.com/strick-j/scimfe/internal/model/user"
	"github.com/strick-j/scimfe/internal/web"
	"go.uber.org/zap"
)

var (
	ErrNotExists = web.NewErrBadRequest("record not found")
	ErrExists    = web.NewErrBadRequest("record already exists")
)

// UserStorage provides user storage
type UserStorage interface {
	// AddUser adds a new user to a storage.
	//
	// Returns a user ID of created user.
	AddUser(ctx context.Context, u user.User) (*user.ID, error)

	// UserByEmail finds user by email
	UserByEmail(ctx context.Context, email string) (*user.User, error)

	// UserByID returns user by ID
	UserByID(ctx context.Context, uid user.ID) (*user.User, error)

	// AllUsers returns all users
	AllUsers(ctx context.Context) (user.Users, error)

	// Exists checks if user with specified email exists
	Exists(email string) (bool, error)
}

type UsersService struct {
	log   *zap.Logger
	store UserStorage
}

func NewUsersService(log *zap.Logger, store UserStorage) *UsersService {
	return &UsersService{
		log:   log.Named("service.users"),
		store: store,
	}
}

// GetAll returns all users
func (s UsersService) GetAll(ctx context.Context) (user.Users, error) {
	return s.store.AllUsers(ctx)
}

// UserByID returns user by id
func (s UsersService) UserByID(ctx context.Context, uid user.ID) (*user.User, error) {
	return s.store.UserByID(ctx, uid)
}

// UserByEmail finds user by email
func (s UsersService) UserByEmail(ctx context.Context, email string) (*user.User, error) {
	usr, err := s.store.UserByEmail(ctx, strings.ToLower(email))
	if err == ErrNotExists {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return usr, nil
}

// AddUser registers a new user
func (s UsersService) AddUser(ctx context.Context, usrReg user.Registration) (*user.User, error) {
	if err := model.Validate(usrReg); err != nil {
		return nil, err
	}

	usrReg.Email = strings.ToLower(usrReg.Email)
	exists, err := s.store.Exists(usrReg.Email)
	if err != nil {
		return nil, fmt.Errorf("can't check if user exists: %w", err)
	}

	if exists {
		return nil, ErrExists
	}

	usr := user.User{Props: usrReg.Props}
	if err = usr.SetPassword(usrReg.Password); err != nil {
		return nil, err
	}

	uid, err := s.store.AddUser(ctx, usr)
	if err != nil {
		return nil, fmt.Errorf("failed to create new user %q: %w", usrReg.Email, err)
	}

	usr.ID = *uid
	return &usr, nil
}
