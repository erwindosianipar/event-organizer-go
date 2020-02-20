package models

import "github.com/jinzhu/gorm"

type Event struct {
	gorm.Model
	ID_User    uint    `gorm:"id_user" json:"id_user,omitempty"`
	Name       string `gorm:"name" json:"name,omitempty"`
	Lokasi     string `gorm:"lokasi" json:"lokasi,omitempty"`
	Event_date string `gorm:"event_date" json:"event_date,omitempty"`
	Kuota      int    `gorm:"kuota" json:"kuota,omitempty"`
	Harga      int    `gorm:"harga" json:"harga,omitempty"`
	Banner     []Banner `gorm:"foreignkey:EventRefer" json:"banner,omitempty"`
}
