package objectMother

import (
	"Booksmart/internal/tests_for_testing/unitTests/serviceTests/testDataBuilder"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type BookModelObjectMother struct {
}

func NewBookModelObjectMother() *BookModelObjectMother {
	return &BookModelObjectMother{}
}

func (bmom *BookModelObjectMother) DefaultBook() *models.BookModel {
	return testDataBuilder.NewBookModelBuilder().Build()
}
