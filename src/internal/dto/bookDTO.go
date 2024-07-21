package dto

type BookParamsDTO struct {
	Title          string
	Author         string
	Publisher      string
	CopiesNumber   int
	Rarity         string
	Genre          string
	PublishingYear int
	Language       string
	AgeLimit       int
	Limit          int // Число элементов на странице
	Offset         int // смещение
}
