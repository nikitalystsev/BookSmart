package models

import (
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type ReaderModelBuilder struct {
	ID          uuid.UUID
	Fio         string
	PhoneNumber string
	Age         uint
	Password    string
	Role        string
}

func NewReaderModelBuilder() *ReaderModelBuilder {
	return &ReaderModelBuilder{
		ID:          uuid.New(),
		Fio:         "default fio",
		PhoneNumber: "00000000000",
		Age:         20,
		Password:    "password00",
		Role:        "Reader",
	}
}

func (builder *ReaderModelBuilder) Build() *models.ReaderModel {
	return &models.ReaderModel{
		ID:          builder.ID,
		Fio:         builder.Fio,
		PhoneNumber: builder.PhoneNumber,
		Age:         builder.Age,
		Password:    builder.Password,
		Role:        builder.Role,
	}
}

func (builder *ReaderModelBuilder) WithID(ID uuid.UUID) *ReaderModelBuilder {
	builder.ID = ID
	return builder
}

func (builder *ReaderModelBuilder) WithFio(fio string) *ReaderModelBuilder {
	builder.Fio = fio
	return builder
}

func (builder *ReaderModelBuilder) WithPhoneNumber(phoneNumber string) *ReaderModelBuilder {
	builder.PhoneNumber = phoneNumber
	return builder
}

func (builder *ReaderModelBuilder) WithAge(age uint) *ReaderModelBuilder {
	builder.Age = age
	return builder
}

func (builder *ReaderModelBuilder) WithPassword(password string) *ReaderModelBuilder {
	builder.Password = password
	return builder
}

func (builder *ReaderModelBuilder) WithRole(role string) *ReaderModelBuilder {
	builder.Role = role
	return builder
}
