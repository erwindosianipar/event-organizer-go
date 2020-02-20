package models

type User struct {
	OrmModel
	Email            string `gorm:"email" json:"email,omitempty"`
	Password         string `gorm:"password" json:"password,omitempty"`
	Name             string `gorm:"name" json:"name,omitempty"`
	Avatar           string `gorm:"avatar" sql:"DEFAULT:'default.png'" json:"avatar,omitempty"`
	Role             string `gorm:"role" sql:"DEFAULT:'user'" json:"role,omitempty"`
	SubmissionStatus string `gorm:"submission_status" sql:"DEFAULT:'not_submit'" json:"submission_status"`
	EventOrganizer   `json:"event_organizer,omitempty"`
}
