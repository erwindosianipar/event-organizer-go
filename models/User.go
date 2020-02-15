package models

import "github.com/jinzhu/gorm"

type EO struct {
	NameEo     string `gorm:"name_eo" sql:"DEFAULT:null" json:"name_eo,omitempty"`
	KTPNumber  string `gorm:"ktp_number" sql:"DEFAULT:null" json:"ktp_number,omitempty"`
	KTPPhoto   string `gorm:"ktp_photo" sql:"DEFAULT:null" json:"ktp_photo,omitempty"`
	SIUPNumber string `gorm:"siup_number" sql:"DEFAULT:null" json:"siup_number,omitempty"`
	IsVerify   bool `gorm:"is_verify" sql:"DEFAULT:false" json:"is_verify,omitempty"`
}

type User struct {
	gorm.Model
	Email       string `gorm:"email" json:"email"`
	Password    string `gorm:"password" json:"password"`
	Name        string `gorm:"name" json:"name"`
	Avatar      string `gorm:"avatar" json:"avatar,omitempty"`
	Role        string `gorm:"role" sql:"DEFAULT:'user'" json:"role,omitempty"`
	EventOrganizer EO `json:"event_organizer,omitempty"`
}
