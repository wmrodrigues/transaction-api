package service

import (
	"context"
	"errors"
	"fmt"
	"transaction-api/internal/common"
	"transaction-api/internal/database"
	"transaction-api/internal/domain"
	"transaction-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetById(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error)
}

type userService struct {
	userRepository     repository.UserRepository
	accountRepository  repository.AccountRepository
	transactionManager database.Manager
}

func NewUserService(userRepository repository.UserRepository, accountRepository repository.AccountRepository, transactionManager database.Manager) UserService {
	return &userService{userRepository: userRepository, accountRepository: accountRepository, transactionManager: transactionManager}
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
	return us.transactionManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := user.Validate(); err != nil {
			return fmt.Errorf("error validating user data: %w", err)
		}
		_, err := us.userRepository.GetByEmail(ctx, user.Email)
		if err == nil {
			return errors.New("user already exists")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error checking existing user: %w", err)
		}
		password, err := common.HashPassword(user.Password)
		if err != nil {
			return fmt.Errorf("error hashing user password: %w", err)
		}
		user.Password = password
		user.Active = true
		user.ID = uuid.New().String()
		if err := us.userRepository.Create(ctx, user); err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
		// here we use the SupportedCurrencies map to create accounts based on the registered currencies
		for currency := range domain.SupportedCurrencies {
			account := &domain.Account{
				ID:       uuid.New().String(),
				UserID:   user.ID,
				Currency: string(currency),
			}
			if err := us.accountRepository.Create(ctx, account); err != nil {
				return fmt.Errorf("error creating %s account for user %s: %w", currency, user.ID, err)
			}
		}
		return nil
	})
}

func (us *userService) ValidateCredentials(ctx context.Context, email, password string) (*domain.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	user, err := us.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user by email, %w", err)
	}
	if !common.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid user credentials")
	}
	return user, nil
}
