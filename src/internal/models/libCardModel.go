package models

import (
	"github.com/google/uuid"
	"time"
)

type LibCardModel struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	ReaderID     uuid.UUID `gorm:"type:uuid" json:"reader_id"`
	LibCardNum   string    `gorm:"type:varchar(13)" json:"lib_card_num"`
	Validity     int       `json:"validity"` // Срок действия
	IssueDate    time.Time `json:"issue_date"`
	ActionStatus bool      `json:"action_status"`
}
