package models

import "time"

type OrmModel struct {
	ID        uint `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}
