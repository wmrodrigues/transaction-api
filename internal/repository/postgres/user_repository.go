package postgres

import (
	"context"
	"fmt"
	"transaction-api/internal/domain"
	"transaction-api/internal/repository"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

func (u UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}
	return &user, nil
}

func (u UserRepository) Create(ctx context.Context, user *domain.User) error {
	err := u.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return err
}

func (u UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}
	return &user, nil
}

func toDomain(model UserModel) *domain.User {
	return &domain.User{
		ID:       model.ID,
		Name:     model.Name,
		Email:    model.Email,
		Password: model.Password,
		Active:   model.Active,
	}
}
func toModel(user *domain.User) UserModel {
	return UserModel{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Active:   user.Active,
	}
}
