package dto

import (
	"errors"
	"transaction-api/internal/domain"
)

type UserRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func (u *UserRequest) Validate() error {
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

func (u *UserRequest) ToDomain() domain.User {
	return domain.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

type UserResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

func UserToResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Active: user.Active,
	}
}
