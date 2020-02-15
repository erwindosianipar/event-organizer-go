package models

type ResponseWrapper struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
