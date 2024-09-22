package testDataBuilder

import (
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type BookModelBuilder struct {
	ID             uuid.UUID
	Title          string
	Author         string
	Publisher      string
	CopiesNumber   uint
	Rarity         string
	Genre          string
	PublishingYear uint
	Language       string
	AgeLimit       uint
}

func NewBookModelBuilder() *BookModelBuilder {
	return &BookModelBuilder{
		ID:             uuid.New(),
		Title:          "default title",
		Author:         "default author",
		Publisher:      "default publisher",
		CopiesNumber:   10,
		Rarity:         "Common",
		Genre:          "default genre",
		PublishingYear: 2021,
		Language:       "english",
		AgeLimit:       0,
	}
}

func (builder *BookModelBuilder) Build() *models.BookModel {
	return &models.BookModel{
		ID:             builder.ID,
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

func (builder *BookModelBuilder) WithID(ID uuid.UUID) *BookModelBuilder {
	builder.ID = ID
	return builder
}

func (builder *BookModelBuilder) WithTitle(title string) *BookModelBuilder {
	builder.Title = title
	return builder
}

func (builder *BookModelBuilder) WithAuthor(author string) *BookModelBuilder {
	builder.Author = author
	return builder
}

func (builder *BookModelBuilder) WithPublisher(publisher string) *BookModelBuilder {
	builder.Publisher = publisher
	return builder
}

func (builder *BookModelBuilder) WithCopiesNumber(copiesNumber uint) *BookModelBuilder {
	builder.CopiesNumber = copiesNumber
	return builder
}

func (builder *BookModelBuilder) WithRarity(rarity string) *BookModelBuilder {
	builder.Rarity = rarity
	return builder
}

func (builder *BookModelBuilder) WithGenre(genre string) *BookModelBuilder {
	builder.Genre = genre
	return builder
}

func (builder *BookModelBuilder) WithPublishingYear(publishingYear uint) *BookModelBuilder {
	builder.PublishingYear = publishingYear
	return builder
}

func (builder *BookModelBuilder) WithLanguage(language string) *BookModelBuilder {
	builder.Language = language
	return builder
}

func (builder *BookModelBuilder) WithAgeLimit(ageLimit uint) *BookModelBuilder {
	builder.AgeLimit = ageLimit
	return builder
}
