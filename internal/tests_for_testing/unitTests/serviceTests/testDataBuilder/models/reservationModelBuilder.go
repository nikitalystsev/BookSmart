package models

import (
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"time"
)

type ReservationModelBuilder struct {
	ID         uuid.UUID
	ReaderID   uuid.UUID
	BookID     uuid.UUID
	IssueDate  time.Time
	ReturnDate time.Time
	State      string
}

func NewReservationModelBuilder() *ReservationModelBuilder {
	return &ReservationModelBuilder{
		ID:         uuid.New(),
		ReaderID:   uuid.New(),
		BookID:     uuid.New(),
		IssueDate:  time.Now(),
		ReturnDate: time.Now().AddDate(0, 0, 14),
		State:      impl.ReservationIssued,
	}
}

func (builder *ReservationModelBuilder) Build() *models.ReservationModel {
	return &models.ReservationModel{
		ID:         builder.ID,
		ReaderID:   builder.ReaderID,
		BookID:     builder.BookID,
		IssueDate:  builder.IssueDate,
		ReturnDate: builder.ReturnDate,
		State:      builder.State,
	}
}

func (builder *ReservationModelBuilder) WithID(ID uuid.UUID) *ReservationModelBuilder {
	builder.ID = ID
	return builder
}

func (builder *ReservationModelBuilder) WithReaderID(readerID uuid.UUID) *ReservationModelBuilder {
	builder.ReaderID = readerID
	return builder
}

func (builder *ReservationModelBuilder) WithBookID(bookID uuid.UUID) *ReservationModelBuilder {
	builder.BookID = bookID
	return builder
}

func (builder *ReservationModelBuilder) WithIssueDate(issueDate time.Time) *ReservationModelBuilder {
	builder.IssueDate = issueDate
	return builder
}

func (builder *ReservationModelBuilder) WithReturnDate(returnDate time.Time) *ReservationModelBuilder {
	builder.IssueDate = returnDate
	return builder
}

func (builder *ReservationModelBuilder) WithState(state string) *ReservationModelBuilder {
	builder.State = state
	return builder
}
