package models

import "github.com/jinzhu/gorm"

type Banner struct {
	gorm.Model
	ID_Event     uint    `gorm:"id_event" json:"id_event,omitempty"`
	Banner_Foto string `gorm:"banner_foto" json:"banner_foto,omitempty"`
}
