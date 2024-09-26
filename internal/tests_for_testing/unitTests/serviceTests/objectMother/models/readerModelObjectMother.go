package models

import (
	tdbmodels "Booksmart/internal/tests_for_testing/unitTests/serviceTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type ReaderModelObjectMother struct {
}

func NewReaderModelObjectMother() *ReaderModelObjectMother {
	return &ReaderModelObjectMother{}
}

func (rmom *ReaderModelObjectMother) DefaultReader() *models.ReaderModel {
	return tdbmodels.NewReaderModelBuilder().Build()
}
