package models

import "github.com/google/uuid"

type ReaderModel struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Fio         string    `json:"fio"`
	PhoneNumber string    `json:"phone_number"`
	Age         uint      `json:"age"`
	Password    string    `json:"password"`
}
