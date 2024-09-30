package models

import (
	"github.com/nikitalystsev/BookSmart-services/core/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
)

type ReaderModelObjectMother struct {
}

func NewReaderModelObjectMother() *ReaderModelObjectMother {
	return &ReaderModelObjectMother{}
}

func (rmom *ReaderModelObjectMother) DefaultReader() *models.ReaderModel {
	return tdbmodels.NewReaderModelBuilder().Build()
}
