package models

import (
	models2 "Booksmart/internal/tests_for_testing/unitTests/serviceTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type BookModelObjectMother struct {
}

func NewBookModelObjectMother() *BookModelObjectMother {
	return &BookModelObjectMother{}
}

func (bmom *BookModelObjectMother) DefaultBook() *models.BookModel {
	return models2.NewBookModelBuilder().Build()
}
