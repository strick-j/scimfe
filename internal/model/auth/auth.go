package auth

import "github.com/strick-j/scimfe/internal/model/user"

type Credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Remember bool   `json:"remember"`
}

type LoginResult struct {
	Token   Token      `json:"token"`
	User    *user.User `json:"user"`
	Session *Session   `json:"session"`
}
