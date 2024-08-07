package models

import (
	"github.com/google/uuid"
	"time"
)

type LibCardModel struct {
	ID           uuid.UUID `json:"id" db:"id" bson:"_id"`
	ReaderID     uuid.UUID `json:"reader_id" db:"reader_id" gorm:"type:uuid" bson:"reader_id"`
	LibCardNum   string    `json:"lib_card_num" db:"lib_card_num" gorm:"type:varchar(13)" bson:"lib_card_num"`
	Validity     int       `json:"validity" db:"validity" bson:"validity"`
	IssueDate    time.Time `json:"issue_date" db:"issue_date" bson:"issue_date"`
	ActionStatus bool      `json:"action_status" db:"action_status" bson:"action_status"`
}
