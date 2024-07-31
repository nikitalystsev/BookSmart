package models

import (
	"github.com/google/uuid"
	"time"
)

type ReservationModel struct {
	ID         uuid.UUID `json:"id" db:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	ReaderID   uuid.UUID `json:"reader_id" db:"reader_id" gorm:"type:uuid"`
	BookID     uuid.UUID `json:"book_id" db:"book_id" gorm:"type:uuid"`
	IssueDate  time.Time `json:"issue_date" db:"issue_date"`
	ReturnDate time.Time `json:"return_date" db:"return_date"`
	State      string    `json:"state" db:"state"`
}
