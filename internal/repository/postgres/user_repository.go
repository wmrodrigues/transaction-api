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
	var model UserModel
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}
	user := model.toDomain()
	return &user, nil
}

func (u UserRepository) Create(ctx context.Context, user *domain.User) error {
	model := UserToModel(user)
	err := u.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return err
}

func (u UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}
	user := model.toDomain()
	return &user, nil
}
