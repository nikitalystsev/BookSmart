package dto

import (
	tdbdto "Booksmart/internal/tests_for_testing/unitTests/serviceTests/testDataBuilder/dto"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
)

type ReaderSignInDTOObjectMother struct {
}

func NewReaderSignInDTOObjectMother() *ReaderSignInDTOObjectMother {
	return &ReaderSignInDTOObjectMother{}
}

func (rdom *ReaderSignInDTOObjectMother) DefaultReaderSignInDTO() *dto.ReaderSignInDTO {
	return tdbdto.NewReaderSignInDTOBuilder().Build()
}
