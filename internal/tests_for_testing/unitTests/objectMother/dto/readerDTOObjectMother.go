package dto

import (
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	tdbdto "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/dto"
)

type ReaderSignInDTOObjectMother struct {
}

func NewReaderSignInDTOObjectMother() *ReaderSignInDTOObjectMother {
	return &ReaderSignInDTOObjectMother{}
}

func (rdom *ReaderSignInDTOObjectMother) DefaultReaderSignInDTO() *dto.SignInInputDTO {
	return tdbdto.NewReaderSignInDTOBuilder().Build()
}
