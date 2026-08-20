package postgres

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(150);not null" json:"name"`
	Email     string         `gorm:"type:varchar(150);not null" json:"email"`
	Password  string         `gorm:"type:varchar(150);not null" json:"password,omitempty"`
	Active    bool           `gorm:"default:true" json:"active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (UserModel) TableName() string {
	return "users"
}
