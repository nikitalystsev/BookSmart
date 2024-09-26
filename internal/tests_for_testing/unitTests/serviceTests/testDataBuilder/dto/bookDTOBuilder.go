package dto

import "github.com/nikitalystsev/BookSmart-services/core/dto"

type BookParamsDTOBuilder struct {
	Title          string
	Author         string
	Publisher      string
	CopiesNumber   uint
	Rarity         string
	Genre          string
	PublishingYear uint
	Language       string
	AgeLimit       uint
	Limit          uint
	Offset         int
}

func NewBookParamsDTOBuilder() *BookParamsDTOBuilder {
	return &BookParamsDTOBuilder{
		Title:          "default title",
		Author:         "default author",
		Publisher:      "default publisher",
		CopiesNumber:   10,
		Rarity:         "Common",
		Genre:          "default genre",
		PublishingYear: 2021,
		Language:       "english",
		AgeLimit:       0,
		Limit:          1,
		Offset:         0,
	}
}

func (builder *BookParamsDTOBuilder) Build() *dto.BookParamsDTO {
	return &dto.BookParamsDTO{
		Title:          builder.Title,
		Author:         builder.Author,
		Publisher:      builder.Publisher,
		CopiesNumber:   builder.CopiesNumber,
		Rarity:         builder.Rarity,
		Genre:          builder.Genre,
		PublishingYear: builder.PublishingYear,
		Language:       builder.Language,
		AgeLimit:       builder.AgeLimit,
	}
}

func (builder *BookParamsDTOBuilder) WithTitle(title string) *BookParamsDTOBuilder {
	builder.Title = title
	return builder
}

func (builder *BookParamsDTOBuilder) WithAuthor(author string) *BookParamsDTOBuilder {
	builder.Author = author
	return builder
}

func (builder *BookParamsDTOBuilder) WithPublisher(publisher string) *BookParamsDTOBuilder {
	builder.Publisher = publisher
	return builder
}

func (builder *BookParamsDTOBuilder) WithCopiesNumber(copiesNumber uint) *BookParamsDTOBuilder {
	builder.CopiesNumber = copiesNumber
	return builder
}

func (builder *BookParamsDTOBuilder) WithRarity(rarity string) *BookParamsDTOBuilder {
	builder.Rarity = rarity
	return builder
}

func (builder *BookParamsDTOBuilder) WithGenre(genre string) *BookParamsDTOBuilder {
	builder.Genre = genre
	return builder
}

func (builder *BookParamsDTOBuilder) WithPublishingYear(publishingYear uint) *BookParamsDTOBuilder {
	builder.PublishingYear = publishingYear
	return builder
}

func (builder *BookParamsDTOBuilder) WithLanguage(language string) *BookParamsDTOBuilder {
	builder.Language = language
	return builder
}

func (builder *BookParamsDTOBuilder) WithAgeLimit(ageLimit uint) *BookParamsDTOBuilder {
	builder.AgeLimit = ageLimit
	return builder
}

func (builder *BookParamsDTOBuilder) WithLimit(limit uint) *BookParamsDTOBuilder {
	builder.Limit = limit
	return builder
}

func (builder *BookParamsDTOBuilder) WithOffset(offset int) *BookParamsDTOBuilder {
	builder.Offset = offset
	return builder
}
