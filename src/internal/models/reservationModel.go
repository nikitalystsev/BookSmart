package models

import (
	"github.com/google/uuid"
	"time"
)

type ReservationModel struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	ReaderID   uuid.UUID `gorm:"type:uuid" json:"reader_id"`
	BookID     uuid.UUID `gorm:"type:uuid" json:"book_id"`
	IssueDate  time.Time `json:"issue_date"`
	ReturnDate time.Time `json:"return_date"`
	State      string    `json:"state"`
}
