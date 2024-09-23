package dto

import (
	dto2 "Booksmart/internal/tests_for_testing/unitTests/serviceTests/testDataBuilder/dto"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
)

type BookParamsDTOObjectMother struct {
}

func NewBookParamsDTOObjectMother() *BookParamsDTOObjectMother {
	return &BookParamsDTOObjectMother{}
}

func (bmom *BookParamsDTOObjectMother) DefaultBookParams() *dto.BookParamsDTO {
	return dto2.NewBookParamsDTOBuilder().Build()
}
