package models

type UserNoPassword struct {
	OrmModel
	Email       string `gorm:"email" json:"email,omitempty"`
	Name        string `gorm:"name" json:"name,omitempty"`
	Avatar      string `gorm:"avatar" sql:"'default.png'" json:"avatar,omitempty"`
	Role        string `gorm:"role" sql:"DEFAULT:'user'" json:"role,omitempty"`
	SubmissionStatus string `gorm:"submission_status" sql:"DEFAULT:'not_submit'" json:"submission_status"`
	EventOrganizer `json:"event_organizer,omitempty"`
}
