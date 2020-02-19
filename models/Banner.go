package models

import "github.com/jinzhu/gorm"

type Banner struct {
	gorm.Model
	Banner_Foto string `gorm:"banner_foto" json:"banner_foto,omitempty"`
	EventRefer uint
}
