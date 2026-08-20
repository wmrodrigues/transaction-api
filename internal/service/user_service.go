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
	GetById(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	ValidateCredentials(ctx context.Context, email, password string) error
}

type userService struct {
	userRepository        repository.UserRepository
	transactionRepository repository.TransactionRepository
}

func NewUserService(userRepository repository.UserRepository, transactionRepository repository.TransactionRepository) UserService {
	return &userService{userRepository: userRepository, transactionRepository: transactionRepository}
}

func (us *userService) GetById(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("user id cannot be empty and must be a valid UUID")
	}
	user, err := us.userRepository.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (us *userService) Create(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("error validating user data: %w", err)
	}
	user.Active = true
	_, err := us.userRepository.GetByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("user already exists")
	}
	user.ID = uuid.New().String()
	err = us.userRepository.Create(ctx, user)
	// create a transaction record with zero value for the new user
	_ = us.transactionRepository.Create(ctx, &domain.Transaction{
		UserID:  user.ID,
		Amount:  0,
		Balance: 0,
	})
	return err
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
