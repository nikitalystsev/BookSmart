package models

import "github.com/google/uuid"

type BookModel struct {
	ID             uuid.UUID `json:"id" db:"id" bson:"_id"`
	Title          string    `json:"title" db:"title" bson:"title"`
	Author         string    `json:"author" db:"author" bson:"author"`
	Publisher      string    `json:"publisher" db:"publisher" bson:"publisher"`
	CopiesNumber   uint      `json:"copies_number" db:"copies_number" bson:"copies_number"`
	Rarity         string    `json:"rarity" db:"rarity" bson:"rarity"`
	Genre          string    `json:"genre" db:"genre" bson:"genre"`
	PublishingYear uint      `json:"publishing_year" db:"publishing_year" bson:"publishing_year"`
	Language       string    `json:"language" db:"language" bson:"language"`
	AgeLimit       uint      `json:"age_limit" db:"age_limit" bson:"age_limit"`
}
