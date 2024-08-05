package models

import "github.com/google/uuid"

type BookModel struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	Author         string    `json:"author" db:"author"`
	Publisher      string    `json:"publisher" db:"publisher"`
	CopiesNumber   uint      `json:"copies_number" db:"copies_number"`
	Rarity         string    `json:"rarity" db:"rarity"`
	Genre          string    `json:"genre" db:"genre"`
	PublishingYear uint      `json:"publishing_year" db:"publishing_year"`
	Language       string    `json:"language" db:"language"`
	AgeLimit       uint      `json:"age_limit" db:"age_limit"`
}
