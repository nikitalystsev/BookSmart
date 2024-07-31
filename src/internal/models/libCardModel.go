package models

import (
	"github.com/google/uuid"
	"time"
)

type LibCardModel struct {
	ID           uuid.UUID `json:"id" db:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	ReaderID     uuid.UUID `json:"reader_id" db:"reader_id" gorm:"type:uuid"`
	LibCardNum   string    `json:"lib_card_num" db:"lib_card_num" gorm:"type:varchar(13)"`
	Validity     int       `json:"validity" db:"validity"`
	IssueDate    time.Time `json:"issue_date" db:"issue_date"`
	ActionStatus bool      `json:"action_status" db:"action_status"`
}
