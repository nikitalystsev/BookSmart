package models

import (
	"github.com/nikitalystsev/BookSmart-services/core/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/serviceTests/testDataBuilder/models"
)

type BookModelObjectMother struct {
}

func NewBookModelObjectMother() *BookModelObjectMother {
	return &BookModelObjectMother{}
}

func (bmom *BookModelObjectMother) DefaultBook() *models.BookModel {
	return tdbmodels.NewBookModelBuilder().Build()
}
