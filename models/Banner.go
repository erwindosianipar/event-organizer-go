package models

import "github.com/jinzhu/gorm"

type Banner struct {
	gorm.Model
	ID_User     int    `gorm:"id_user" json:"id_user,omitempty"`
	Banner_Foto string `gorm:"banner_foto" json:"banner_foto,omitempty"`
}
