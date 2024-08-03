package dto

type BookParamsDTO struct {
	Title          string
	Author         string
	Publisher      string
	CopiesNumber   uint
	Rarity         string
	Genre          string
	PublishingYear uint
	Language       string
	AgeLimit       uint
	Limit          uint // Число элементов на странице
	Offset         int  // смещение
}
