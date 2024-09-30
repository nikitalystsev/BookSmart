package dto

import "github.com/nikitalystsev/BookSmart-services/core/dto"

type ReaderSignInDTOBuilder struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func NewReaderSignInDTOBuilder() *ReaderSignInDTOBuilder {
	return &ReaderSignInDTOBuilder{
		PhoneNumber: "00000000000",
		Password:    "password00",
	}
}

func (builder *ReaderSignInDTOBuilder) Build() *dto.SignInInputDTO {
	return &dto.SignInInputDTO{
		PhoneNumber: builder.PhoneNumber,
		Password:    builder.Password,
	}
}

func (builder *ReaderSignInDTOBuilder) WithPhoneNumber(phoneNumber string) *ReaderSignInDTOBuilder {
	builder.PhoneNumber = phoneNumber
	return builder
}

func (builder *ReaderSignInDTOBuilder) WithPassword(password string) *ReaderSignInDTOBuilder {
	builder.Password = password
	return builder
}
