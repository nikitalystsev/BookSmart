package models

import "github.com/google/uuid"

type ReaderModel struct {
	ID          uuid.UUID `json:"id" db:"id" bson:"_id"`
	Fio         string    `json:"fio" db:"fio" bson:"fio"`
	PhoneNumber string    `json:"phone_number" db:"phone_number" bson:"phone_number"`
	Age         uint      `json:"age" db:"age" bson:"age"`
	Password    string    `json:"password" db:"password" bson:"password"`
	Role        string    `json:"role" db:"role" bson:"role"`
}
