package postgres

import (
	"time"
	"transaction-api/internal/domain"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"password,omitempty"`
	Active    bool           `gorm:"default:true" json:"active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (u *UserModel) TableName() string {
	return "users"
}

func (u *UserModel) toDomain() domain.User {
	return domain.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Active:   u.Active,
	}
}

func UserToModel(user *domain.User) UserModel {
	return UserModel{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Active:   user.Active,
	}
}
