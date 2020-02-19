package models

type ResponseListEvent struct {
	ID     uint   `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Lokasi string `json:"lokasi,omitempty"`
}
