package integrationTests

import (
	"BookSmart-services/core/dto"
	"BookSmart-services/core/models"
	"BookSmart-services/errs"
	"context"
	"errors"
	"github.com/google/uuid"
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

	_, err := s.readerService.SignIn(context.Background(), readerDTO)
	s.NoError(err)
}

func (s *IntegrationTestSuite) TestReader_SignIn_Error() {

	readerDTO := &dto.ReaderSignInDTO{
		PhoneNumber: "79314562376",
		Password:    "hjghhgdgfs",
	}

	_, err := s.readerService.SignIn(context.Background(), readerDTO)
	s.Error(err)
	s.Equal(errors.New("wrong password"), err)
}
