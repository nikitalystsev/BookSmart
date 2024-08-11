package dto

type BookParamsDTO struct {
	Title          string `json:"title" bson:"title"`
	Author         string `json:"author" bson:"author"`
	Publisher      string `json:"publisher" bson:"publisher"`
	CopiesNumber   uint   `json:"copies_number" bson:"copies_number"`
	Rarity         string `json:"rarity" bson:"rarity"`
	Genre          string `json:"genre" bson:"genre"`
	PublishingYear uint   `json:"publishing_year" bson:"publishing_year"`
	Language       string `json:"language" bson:"language"`
	AgeLimit       uint   `json:"age_limit" bson:"age_limit"`
	Limit          uint   `json:"limit" bson:"limit"`   // Число элементов на странице
	Offset         int    `json:"offset" bson:"offset"` // Смещение
}
