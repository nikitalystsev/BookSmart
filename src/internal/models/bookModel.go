package models

import "github.com/google/uuid"

type BookModel struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	Publisher      string    `json:"publisher"`
	CopiesNumber   int       `json:"copies_number"`
	Rarity         string    `json:"rarity"`
	Genre          string    `json:"genre"`
	PublishingYear int       `json:"publishing_year"`
	Language       string    `json:"language"`
	AgeLimit       int       `json:"age_limit"`
}
