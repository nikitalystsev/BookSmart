package integrationTests

//func (s *IntegrationTestSuite) TestReservation_Create_Success() {
//	readerDTO := &dto.SignInInputDTO{
//		PhoneNumber: "79314562376",
//		Password:    "sdgdgsgsgd",
//	}
//	_, err := s.readerService.SignIn(context.Background(), readerDTO.PhoneNumber, readerDTO.Password)
//	s.NoError(err)
//
//	params := &dto.BookParamsDTO{
//		Title:  "Harry Potter and the Order of the Phoenix",
//		Limit:  1,
//		Offset: 0,
//	}
//
//	readerID, _ := uuid.Parse("75919792-c2d9-4685-92b2-e2a80b2ed5be")
//
//	var books []*models.BookModel
//	books, err = s.bookService.GetByParams(context.Background(), params)
//	s.NoError(err)
//	s.Len(books, 1)
//
//	err = s.reservationService.Create(context.Background(), readerID, books[0].ID)
//	s.NoError(err)
//
//	expectedReservations, err := s.reservationService.GetByReaderID(context.Background(), readerID, impl.ReservationsPageLimit, 0)
//	s.NoError(err)
//	s.Len(expectedReservations, 3)
//
//	var expectedReservation *models.ReservationModel
//	for _, _expectedReservation := range expectedReservations {
//		if _expectedReservation.BookID == books[0].ID {
//			expectedReservation = _expectedReservation
//		}
//	}
//	s.Equal(expectedReservation.BookID, books[0].ID)
//	s.Equal(expectedReservation.ReaderID, readerID)
//
//}
//
//func (s *IntegrationTestSuite) TestReservation_Create_Error() {
//	readerDTO := &dto.SignInInputDTO{
//		PhoneNumber: "76867456521",
//		Password:    "hghhfnnbdd",
//	}
//	_, err := s.readerService.SignIn(context.Background(), readerDTO.PhoneNumber, readerDTO.Password)
//	s.NoError(err)
//
//	params := &dto.BookParamsDTO{
//		Title:  "Harry Potter and the Order of the Phoenix",
//		Limit:  1,
//		Offset: 0,
//	}
//
//	readerID, _ := uuid.Parse("3885b2d3-ef6e-4f62-8f86-d1454d108207")
//
//	var books []*models.BookModel
//	books, err = s.bookService.GetByParams(context.Background(), params)
//	s.NoError(err)
//	s.Len(books, 1)
//
//	err = s.reservationService.Create(context.Background(), readerID, books[0].ID)
//	s.Error(err)
//	s.Equal(errs.ErrLibCardIsInvalid, err)
//}
//
//func (s *IntegrationTestSuite) TestReservation_Update_Success() {
//	readerDTO := &dto.SignInInputDTO{
//		PhoneNumber: "79314562376",
//		Password:    "sdgdgsgsgd",
//	}
//
//	_, err := s.readerService.SignIn(context.Background(), readerDTO.PhoneNumber, readerDTO.Password)
//	s.NoError(err)
//
//	readerID, err := uuid.Parse("75919792-c2d9-4685-92b2-e2a80b2ed5be")
//	bookID, err := uuid.Parse("43f45552-4a95-4f12-864b-e1d8bfa30b8d")
//
//	reservations, err := s.reservationService.GetByReaderID(context.Background(), readerID, impl.ReservationsPageLimit, 0)
//	s.NoError(err)
//
//	var testReservation *models.ReservationModel
//	for _, reservation := range reservations {
//		if reservation.BookID == bookID {
//			testReservation = reservation
//			break
//		}
//	}
//
//	s.Equal(impl.ReservationIssued, testReservation.State)
//
//	err = s.reservationService.Update(context.Background(), testReservation, 5)
//	s.NoError(err)
//	expectedReservation, err := s.reservationService.GetByID(context.Background(), testReservation.ID)
//	s.NoError(err)
//	s.Equal(impl.ReservationExtended, expectedReservation.State)
//}
//
//func (s *IntegrationTestSuite) TestReservation_Update_Error() {
//	readerDTO := &dto.SignInInputDTO{
//		PhoneNumber: "79314562376",
//		Password:    "sdgdgsgsgd",
//	}
//
//	_, err := s.readerService.SignIn(context.Background(), readerDTO.PhoneNumber, readerDTO.Password)
//	s.NoError(err)
//
//	readerID, err := uuid.Parse("5818061a-662d-45bb-a67c-0d2873038e65")
//	bookID, err := uuid.Parse("b33b30c8-254e-45f2-8314-0b93a6b8c561")
//
//	reservations, err := s.reservationService.GetByReaderID(context.Background(), readerID, impl.ReservationsPageLimit, 0)
//	s.NoError(err)
//
//	var testReservation *models.ReservationModel
//	for _, reservation := range reservations {
//		if reservation.BookID == bookID {
//			testReservation = reservation
//			break
//		}
//	}
//
//	err = s.reservationService.Update(context.Background(), testReservation, 5)
//	s.Error(err)
//	s.Error(errs.ErrRareAndUniqueBookNotExtended, err)
//}
