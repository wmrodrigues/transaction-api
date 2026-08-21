package dto

import (
	"errors"
	"transaction-api/internal/domain"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (u *User) ToDomain() domain.User {
	return domain.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}
