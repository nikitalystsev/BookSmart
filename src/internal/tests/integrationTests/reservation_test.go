package integrationTests

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"context"
	"github.com/google/uuid"
)

func (s *IntegrationTestSuite) TestReservation_Create() {
	reader := &models.ReaderModel{
		ID:          uuid.New(),
		Fio:         "John Doe",
		PhoneNumber: "79214553812",
		Age:         30,
		Password:    "fdgsshhyrs",
	}

	err := s.readerService.SignUp(context.Background(), reader)
	s.NoError(err)

	readerDTO := &dto.ReaderSignInDTO{
		PhoneNumber: "79214553812",
		Password:    "fdgsshhyrs",
	}

	_, err = s.readerService.SignIn(context.Background(), readerDTO)
	s.NoError(err)

	params := &dto.BookParamsDTO{
		Title:  "Harry Potter and the Order of the Phoenix",
		Limit:  1,
		Offset: 0,
	}

	err = s.libCardService.Create(context.Background(), reader.ID)
	s.NoError(err)
	
	var books []*models.BookModel
	books, err = s.bookService.GetByParams(context.Background(), params)
	s.NoError(err)
	s.Len(books, 1)

	err = s.reservationService.Create(context.Background(), reader.ID, books[0].ID)
	s.NoError(err)
}
