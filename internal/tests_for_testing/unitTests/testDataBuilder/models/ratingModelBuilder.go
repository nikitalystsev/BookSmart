package models

import (
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type RatingModelBuilder struct {
	ID       uuid.UUID
	ReaderID uuid.UUID
	BookID   uuid.UUID
	Review   string
	Rating   int
}

func NewRatingModelBuilder() *RatingModelBuilder {
	return &RatingModelBuilder{
		ID:       uuid.New(),
		ReaderID: uuid.New(),
		BookID:   uuid.New(),
		Review:   "ok",
		Rating:   5,
	}
}

func (builder *RatingModelBuilder) Build() *models.RatingModel {
	return &models.RatingModel{
		ID:       builder.ID,
		ReaderID: builder.ReaderID,
		BookID:   builder.BookID,
		Review:   builder.Review,
		Rating:   builder.Rating,
	}
}

func (builder *RatingModelBuilder) WithID(ID uuid.UUID) *RatingModelBuilder {
	builder.ID = ID
	return builder
}

func (builder *RatingModelBuilder) WithReaderID(readerID uuid.UUID) *RatingModelBuilder {
	builder.ReaderID = readerID
	return builder
}

func (builder *RatingModelBuilder) WithBookID(bookID uuid.UUID) *RatingModelBuilder {
	builder.BookID = bookID
	return builder
}

func (builder *RatingModelBuilder) WithReview(review string) *RatingModelBuilder {
	builder.Review = review
	return builder
}

func (builder *RatingModelBuilder) WithRating(rating int) *RatingModelBuilder {
	builder.Rating = rating
	return builder
}
