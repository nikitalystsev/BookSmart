package models

import (
	"github.com/google/uuid"
	"time"
)

type ReservationModel struct {
	ID         uuid.UUID `json:"id" db:"id" bson:"_id"`
	ReaderID   uuid.UUID `json:"reader_id" db:"reader_id" gorm:"type:uuid" bson:"reader_id"`
	BookID     uuid.UUID `json:"book_id" db:"book_id" gorm:"type:uuid" bson:"book_id"`
	IssueDate  time.Time `json:"issue_date" db:"issue_date" bson:"issue_date"`
	ReturnDate time.Time `json:"return_date" db:"return_date" bson:"return_date"`
	State      string    `json:"state" db:"state" bson:"state"`
}
