package unitTests

//func TestReaderService_Register(t *testing.T) {
//
//	type mockBehaviour func(m *mockrepositories.MockIReaderRepo, reader *models.ReaderModel)
//
//	testTable := []struct {
//		name          string
//		mockBehaviour mockBehaviour
//		expectedError error
//	}{
//		{
//			name: "successfully register reader",
//			mockBehaviour: func(m *mockrepositories.MockIReaderRepo, reader *models.ReaderModel) {
//				m.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, nil)
//				m.EXPECT().Create(gomock.Any(), reader).Return(nil)
//			},
//			expectedError: nil,
//		},
//	}
//
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//
//			mockReaderRepo := mockrepositories.NewMockIReaderRepo(ctrl)
//			readerService := implServices.NewReaderService(
//				mockReaderRepo,
//			)
//
//			newReader := &models.ReaderModel{
//				ID:          uuid.New(),
//				Fio:         "Jon Smith",
//				PhoneNumber: "+79313452367",
//				Password:    "password",
//			}
//
//			testCase.mockBehaviour(mockReaderRepo, newReader)
//
//			err := readerService.Register(context.Background(), newReader)
//			assert.Equal(t, testCase.expectedError, err)
//		})
//	}
//}
//
//func TestReaderService_Login(t *testing.T) {
//
//	type mockBehaviour func(m *mockrepositories.MockIReaderRepo, readerModel *models.ReaderModel, readerDTO *dto.ReaderLoginDTO)
//
//	testTable := []struct {
//		name          string
//		mockBehaviour mockBehaviour
//		expectedError error
//	}{
//		{
//			name: "successfully login reader",
//			mockBehaviour: func(m *mockrepositories.MockIReaderRepo, readerModel *models.ReaderModel, readerDTO *dto.ReaderLoginDTO) {
//				m.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(readerModel, nil)
//			},
//			expectedError: nil,
//		},
//	}
//
//	for _, testCase := range testTable {
//		t.Run(testCase.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//
//			mockReaderRepo := mockrepositories.NewMockIReaderRepo(ctrl)
//			readerService := implServices.CreateNewReaderService(
//				mockReaderRepo,
//			)
//
//			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
//
//			newReader := &models.ReaderModel{
//				ID:          uuid.New(),
//				PhoneNumber: "+79313452367",
//				Password:    string(hashedPassword),
//			}
//
//			newReaderDTO := &dto.ReaderLoginDTO{
//				PhoneNumber: "+79313452367",
//				Password:    "password",
//			}
//
//			testCase.mockBehaviour(mockReaderRepo, newReader, newReaderDTO)
//
//			err := readerService.Login(context.Background(), newReaderDTO)
//			assert.NoError(t, err)
//		})
//	}
//}
