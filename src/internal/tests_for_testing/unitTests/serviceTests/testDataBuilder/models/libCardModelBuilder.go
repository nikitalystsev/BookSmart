package models

import (
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"time"
)

type LibCardModelBuilder struct {
	ID           uuid.UUID
	ReaderID     uuid.UUID
	LibCardNum   string
	Validity     int
	IssueDate    time.Time
	ActionStatus bool
}

func NewLibCardModelBuilder() *LibCardModelBuilder {
	return &LibCardModelBuilder{
		ID:           uuid.New(),
		ReaderID:     uuid.New(),
		LibCardNum:   "1234567890123",
		Validity:     impl.LibCardValidityPeriod,
		IssueDate:    time.Now(),
		ActionStatus: true,
	}
}

func (builder *LibCardModelBuilder) Build() *models.LibCardModel {
	return &models.LibCardModel{
		ID:           builder.ID,
		ReaderID:     builder.ReaderID,
		LibCardNum:   builder.LibCardNum,
		Validity:     builder.Validity,
		IssueDate:    builder.IssueDate,
		ActionStatus: builder.ActionStatus,
	}
}

func (builder *LibCardModelBuilder) WithID(ID uuid.UUID) *LibCardModelBuilder {
	builder.ID = ID
	return builder
}

func (builder *LibCardModelBuilder) WithReaderID(ReaderID uuid.UUID) *LibCardModelBuilder {
	builder.ReaderID = ReaderID
	return builder
}

func (builder *LibCardModelBuilder) WithLibCardNum(libCardNum string) *LibCardModelBuilder {
	builder.LibCardNum = libCardNum
	return builder
}

func (builder *LibCardModelBuilder) WithValidity(Validity int) *LibCardModelBuilder {
	builder.Validity = Validity
	return builder
}

func (builder *LibCardModelBuilder) WithIssueDate(IssueDate time.Time) *LibCardModelBuilder {
	builder.IssueDate = IssueDate
	return builder
}

func (builder *LibCardModelBuilder) WithActionStatus(ActionStatus bool) *LibCardModelBuilder {
	builder.ActionStatus = ActionStatus
	return builder
}
