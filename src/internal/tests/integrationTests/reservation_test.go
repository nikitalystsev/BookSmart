package integrationTests

import (
	"BookSmart-services/core/dto"
	"BookSmart-services/core/models"
	"BookSmart-services/errs"
	"context"
	"github.com/google/uuid"
)

func (s *IntegrationTestSuite) TestReservation_Create_Success() {
	readerDTO := &dto.ReaderSignInDTO{
		PhoneNumber: "79314562376",
		Password:    "sdgdgsgsgd",
	}
	_, err := s.readerService.SignIn(context.Background(), readerDTO)
	s.NoError(err)

	params := &dto.BookParamsDTO{
		Title:  "Harry Potter and the Order of the Phoenix",
		Limit:  1,
		Offset: 0,
	}

	readerID, _ := uuid.Parse("75919792-c2d9-4685-92b2-e2a80b2ed5be")

	var books []*models.BookModel
	books, err = s.bookService.GetByParams(context.Background(), params)
	s.NoError(err)
	s.Len(books, 1)

	err = s.reservationService.Create(context.Background(), readerID, books[0].ID)
	s.NoError(err)
}

func (s *IntegrationTestSuite) TestReservation_Create_Error() {
	readerDTO := &dto.ReaderSignInDTO{
		PhoneNumber: "32534523451",
		Password:    "rtjhhhgffr",
	}
	_, err := s.readerService.SignIn(context.Background(), readerDTO)
	s.NoError(err)

	params := &dto.BookParamsDTO{
		Title:  "Harry Potter and the Order of the Phoenix",
		Limit:  1,
		Offset: 0,
	}

	readerID, _ := uuid.Parse("6800b3ee-9810-450e-9ca5-776aa1c6191d")

	var books []*models.BookModel
	books, err = s.bookService.GetByParams(context.Background(), params)
	s.NoError(err)
	s.Len(books, 1)

	err = s.reservationService.Create(context.Background(), readerID, books[0].ID)
	s.Error(err)
	s.Equal(errs.ErrLibCardIsInvalid, err)
}

func (s *IntegrationTestSuite) TestReservation_Update_Success() {
	readerDTO := &dto.ReaderSignInDTO{
		PhoneNumber: "79314562376",
		Password:    "sdgdgsgsgd",
	}

	_, err := s.readerService.SignIn(context.Background(), readerDTO)
	s.NoError(err)

	readerID, err := uuid.Parse("75919792-c2d9-4685-92b2-e2a80b2ed5be")
	bookID, err := uuid.Parse("43f45552-4a95-4f12-864b-e1d8bfa30b8d")

	reservations, err := s.reservationService.GetAllReservationsByReaderID(context.Background(), readerID)
	s.NoError(err)

	var testReservation *models.ReservationModel
	for _, reservation := range reservations {
		if reservation.BookID == bookID {
			testReservation = reservation
			break
		}
	}

	err = s.reservationService.Update(context.Background(), testReservation)
	s.NoError(err)
}

func (s *IntegrationTestSuite) TestReservation_Update_Error() {
	readerDTO := &dto.ReaderSignInDTO{
		PhoneNumber: "79314562376",
		Password:    "sdgdgsgsgd",
	}

	_, err := s.readerService.SignIn(context.Background(), readerDTO)
	s.NoError(err)

	readerID, err := uuid.Parse("75919792-c2d9-4685-92b2-e2a80b2ed5be")
	bookID, err := uuid.Parse("f01107fb-4f7a-4f37-ba1e-6c6012c5203c")

	reservations, err := s.reservationService.GetAllReservationsByReaderID(context.Background(), readerID)
	s.NoError(err)

	var testReservation *models.ReservationModel
	for _, reservation := range reservations {
		if reservation.BookID == bookID {
			testReservation = reservation
			break
		}
	}

	err = s.reservationService.Update(context.Background(), testReservation)
	s.Error(err)
	s.Error(errs.ErrRareAndUniqueBookNotExtended, err)
}
