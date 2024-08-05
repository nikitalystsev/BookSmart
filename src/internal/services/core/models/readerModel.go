package models

import "github.com/google/uuid"

type ReaderModel struct {
	ID          uuid.UUID `json:"id" db:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Fio         string    `json:"fio" db:"fio"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Age         uint      `json:"age" db:"age"`
	Password    string    `json:"password" db:"password"`
	Role        string    `json:"role" db:"role"`
}
