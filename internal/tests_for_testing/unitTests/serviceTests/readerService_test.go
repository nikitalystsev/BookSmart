package serviceTests

//func TestReaderService_SingUp_Success(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	var err error
//
//	// Arrange
//	mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
//	mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
//	readerService := impl.NewReaderService(
//		mockReaderRepo, nil, nil,
//		mockHasher, logging.GetLoggerForTests(),
//		1, 2,
//	)
//	reader := ommodels.NewReaderModelObjectMother().DefaultReader()
//	mockReaderRepo.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
//	mockHasher.EXPECT().Hash(reader.Password).Return("hashed password", nil)
//	mockReaderRepo.EXPECT().Create(gomock.Any(), reader).Return(nil)
//
//	runner.Run(t, "success sing up", func(t provider.T) {
//		// Act
//		err = readerService.SignUp(context.Background(), reader)
//	})
//
//	// Assert
//	assert.Nil(t, err)
//}
//
//func TestReaderService_SingUp_ErrorCheckReaderExistence(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	var err error
//
//	// Arrange
//	mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
//	mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
//	readerService := impl.NewReaderService(
//		mockReaderRepo, nil, nil,
//		mockHasher, logging.GetLoggerForTests(),
//		1, 2,
//	)
//	reader := ommodels.NewReaderModelObjectMother().DefaultReader()
//	mockReaderRepo.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errors.New("database error"))
//
//	runner.Run(t, "error check reader existence", func(t provider.T) {
//		// Act
//		err = readerService.SignUp(context.Background(), reader)
//	})
//
//	// Assert
//	assert.NotNil(t, err)
//	assert.Equal(t, errors.New("database error"), err)
//}
//
//func TestReaderService_SingIn_Success(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	var (
//		err    error
//		tokens *models.Tokens
//	)
//
//	// Arrange
//	mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
//	mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
//	mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
//	readerService := impl.NewReaderService(
//		mockReaderRepo, nil, mockTokenManager,
//		mockHasher, logging.GetLoggerForTests(),
//		1, 2,
//	)
//	readerDTO := omdto.NewReaderSignInDTOObjectMother().DefaultReaderSignInDTO()
//	reader := ommodels.NewReaderModelObjectMother().DefaultReader()
//	mockReaderRepo.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
//	mockHasher.EXPECT().Compare(reader.Password, readerDTO.Password).Return(nil)
//	mockTokenManager.EXPECT().NewJWT(reader.ID, reader.Role, 1).Return("accessToken", nil)
//	mockTokenManager.EXPECT().NewRefreshToken().Return("refreshToken", nil)
//	mockReaderRepo.EXPECT().SaveRefreshToken(gomock.Any(), reader.ID, "refreshToken", 2).Return(nil)
//
//	runner.Run(t, "success sign in", func(t provider.T) {
//		// Act
//		tokens, err = readerService.SignIn(context.Background(), readerDTO.PhoneNumber, readerDTO.Password)
//	})
//
//	// Assert
//	assert.Nil(t, err)
//	assert.NotNil(t, tokens)
//}
