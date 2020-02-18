package models

type EventOrganizer struct {
	NameEo     string `gorm:"name_eo" json:"name_eo,omitempty"`
	KTPNumber  string `gorm:"ktp_number" json:"ktp_number,omitempty"`
	KTPPhoto   string `gorm:"ktp_photo" json:"ktp_photo,omitempty"`
	SIUPNumber string `gorm:"siup_number" json:"siup_number,omitempty"`
	IsVerify   bool   `gorm:"is_verify" json:"is_verify"`
}
