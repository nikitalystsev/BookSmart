package integrationTests

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
)

func (s *IntegrationTestSuite) TestReader_SignUp_Success() {
	reader := &models.ReaderModel{
		ID:          uuid.New(),
		Fio:         "John Doe",
		PhoneNumber: "79214553467",
		Age:         30,
		Password:    "gfdggsshdf",
	}

	err := s.readerService.SignUp(context.Background(), reader)
	s.NoError(err)

	expectedReader, err := s.readerService.GetByPhoneNumber(context.Background(), reader.PhoneNumber)
	s.NoError(err)
	s.Equal(reader, expectedReader)
}

func (s *IntegrationTestSuite) TestReader_SignUp_Error() {
	reader := &models.ReaderModel{
		ID:          uuid.New(),
		Fio:         "John Doe",
		PhoneNumber: "79214553467",
		Age:         30,
		Password:    "sdgdesf",
	}

	err := s.readerService.SignUp(context.Background(), reader)
	s.Error(err)
	s.Equal(errs.ErrInvalidReaderPasswordLen, err)
}

func (s *IntegrationTestSuite) TestReader_SignIn_Success() {

	readerDTO := &dto.ReaderSignInDTO{
		PhoneNumber: "79314562376",
		Password:    "sdgdgsgsgd",
	}

	tokens, err := s.readerService.SignIn(context.Background(), readerDTO.PhoneNumber, readerDTO.Password)
	s.NoError(err)
	s.NotNil(tokens)
}

func (s *IntegrationTestSuite) TestReader_SignIn_Error() {

	readerDTO := &dto.ReaderSignInDTO{
		PhoneNumber: "79314562376",
		Password:    "hjghhgdgfs",
	}

	_, err := s.readerService.SignIn(context.Background(), readerDTO.PhoneNumber, readerDTO.Password)
	s.Error(err)
	s.Equal(errors.New("wrong password"), err)
}
