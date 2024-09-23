package models

import (
	models2 "Booksmart/internal/tests_for_testing/unitTests/serviceTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"time"
)

type LibCardModelObjectMother struct {
}

func NewLibCardModelObjectMother() *LibCardModelObjectMother {
	return &LibCardModelObjectMother{}
}

func (lcmom *LibCardModelObjectMother) DefaultLibCard() *models.LibCardModel {
	return models2.NewLibCardModelBuilder().Build()
}

func (lcmom *LibCardModelObjectMother) ExpiredLibCard() *models.LibCardModel {
	return models2.NewLibCardModelBuilder().
		WithIssueDate(time.Now().AddDate(0, 0, -370)).
		WithActionStatus(false).
		Build()
}
