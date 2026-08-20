package service

import (
	"context"
	"errors"
	"fmt"
	"transaction-api/internal/domain"
	"transaction-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	ValidateCredentials(ctx context.Context, email, password string) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (us *userService) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("user id cannot be empty and must be a valid UUID")
	}
	user, err := us.userRepository.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (us *userService) CreateUser(ctx context.Context, user *domain.User) error {
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	user.Active = true
	_, err := us.userRepository.GetByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("user already exists")
	}
	user.ID = uuid.New().String()
	return us.userRepository.Create(ctx, user)
}

func (us *userService) ValidateCredentials(ctx context.Context, email, password string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if password == "" {
		return errors.New("password is required")
	}
	user, err := us.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("error getting user by email, %w", err)
	}
	if user.Password != password {
		return errors.New("invalid user credentials")
	}
	return nil
}
