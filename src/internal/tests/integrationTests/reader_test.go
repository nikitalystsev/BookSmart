package integrationTests

//func (s *IntegrationTestSuite) TestReader_SignUp() {
//	reader := &models.ReaderModel{
//		ID:          uuid.New(),
//		Fio:         "John Doe",
//		PhoneNumber: "79214553467",
//		Age:         30,
//		Password:    "gfdggsshdf",
//	}
//
//	err := s.readerService.SignUp(context.Background(), reader)
//	s.NoError(err)
//}
//
//func (s *IntegrationTestSuite) TestReader_SignIn() {
//	reader := &models.ReaderModel{
//		ID:          uuid.New(),
//		Fio:         "John Doe",
//		PhoneNumber: "79215672398",
//		Age:         30,
//		Password:    "dgfdrshyde",
//	}
//
//	err := s.readerService.SignUp(context.Background(), reader)
//	s.NoError(err)
//
//	readerDTO := &dto.ReaderSignInDTO{
//		PhoneNumber: "79215672398",
//		Password:    "dgfdrshyde",
//	}
//
//	_, err = s.readerService.SignIn(context.Background(), readerDTO)
//
//	s.NoError(err)
//}
