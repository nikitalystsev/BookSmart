package serviceTests

/*
Лондонский вариант
*/

//func TestLibCardService_Create_Success(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	var err error
//	// Arrange
//	mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
//	libCardService := impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
//	readerID := uuid.New()
//	mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errs.ErrLibCardDoesNotExists)
//	mockLibCardRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
//
//	runner.Run(t, "success create libCard", func(t provider.T) {
//
//		// Act
//		err = libCardService.Create(context.Background(), readerID)
//	})
//
//	// Assert
//	assert.Nil(t, err)
//}
//
//func TestLibCardService_Create_ErrorLibCardAlreadyExists(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	var err error
//
//	// Arrange
//	mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
//	libCardService := impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
//	readerID := uuid.New()
//	libCard := tdbmodels.NewLibCardModelBuilder().WithReaderID(readerID).Build()
//	mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(libCard, nil)
//
//	runner.Run(t, "error libCard already exist", func(t provider.T) {
//
//		// Act
//		err = libCardService.Create(context.Background(), readerID)
//	})
//
//	// Assert
//	assert.NotNil(t, err)
//	assert.Equal(t, errs.ErrLibCardAlreadyExist, err)
//}
//
//func TestLibCardService_Update_Success(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	var err error
//
//	// Arrange
//	mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
//	libCardService := impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
//	libCard := ommodels.NewLibCardModelObjectMother().ExpiredLibCard()
//	mockLibCardRepo.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(libCard, nil)
//	mockLibCardRepo.EXPECT().Update(gomock.Any(), libCard).Return(nil)
//
//	runner.Run(t, "success update libCard", func(t provider.T) {
//
//		// Act
//		err = libCardService.Update(context.Background(), libCard)
//	})
//
//	// Assert
//	assert.Nil(t, err)
//}
//
//func TestLibCardService_Update_ErrorCheckLibCardExistence(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	var err error
//
//	// Arrange
//	mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
//	libCardService := impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
//	libCard := ommodels.NewLibCardModelObjectMother().DefaultLibCard()
//	mockLibCardRepo.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(nil, errors.New("database error"))
//
//	runner.Run(t, "error check libCard existence", func(t provider.T) {
//
//		// Act
//		err = libCardService.Update(context.Background(), libCard)
//	})
//
//	// Assert
//	assert.NotNil(t, err)
//	assert.Equal(t, errors.New("database error"), err)
//}
//
//func TestLibCardService_GetByReaderID_Success(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	var (
//		err         error
//		findLibCard *models.LibCardModel
//	)
//
//	// Arrange
//	mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
//	libCardService := impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
//	readerID := uuid.New()
//	libCard := tdbmodels.NewLibCardModelBuilder().WithReaderID(readerID).Build()
//	mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(libCard, nil)
//
//	runner.Run(t, "success get libCard by readerID", func(t provider.T) {
//
//		// Act
//		findLibCard, err = libCardService.GetByReaderID(context.Background(), readerID)
//	})
//
//	// Assert
//	assert.Nil(t, err)
//	assert.Equal(t, libCard, findLibCard)
//}
//
//func TestLibCardService_GetByReaderID_ErrorLibCardDoesNotExists(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	var (
//		err         error
//		findLibCard *models.LibCardModel
//	)
//
//	// Arrange
//	mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
//	libCardService := impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
//	readerID := uuid.New()
//	mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errs.ErrLibCardDoesNotExists)
//
//	runner.Run(t, "error libCard does not exist", func(t provider.T) {
//		// Act
//		findLibCard, err = libCardService.GetByReaderID(context.Background(), readerID)
//	})
//
//	// Assert
//	assert.NotNil(t, err)
//	assert.Equal(t, errs.ErrLibCardDoesNotExists, err)
//	assert.Nil(t, findLibCard)
//}
//
//func TestLibCardService_Create_Success_Classic(t *testing.T) {
//	ctx := context.Background()
//	var err error
//
//	// Arrange
//	container, err := getContainerForClassicUnitTests()
//	if err != nil {
//		t.Fatal(err)
//	}
//	db, err := applyMigrations(container)
//	if err != nil {
//		t.Fatal(err)
//	}
//	reader := ommodels.NewReaderModelObjectMother().DefaultReader()
//
//	_, _ = db.ExecContext(
//		ctx, `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
//		reader.ID,
//		reader.Fio,
//		reader.PhoneNumber,
//		reader.Age,
//		reader.Password,
//		reader.Role,
//	)
//	libCardRepo := implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())
//	libCardService := impl.NewLibCardService(libCardRepo, logging.GetLoggerForTests())
//	defer func(db *sqlx.DB) {
//		if err = db.Close(); err != nil {
//			t.Fatalf("failed to close database connection: %v\n", err)
//		}
//	}(db)
//
//	defer func() {
//		if err = container.Terminate(ctx); err != nil {
//			t.Fatalf("failed to terminate container: %v\n", err)
//		}
//	}()
//
//	runner.Run(t, "success create libCard", func(t provider.T) {
//
//		// Act
//		err = libCardService.Create(context.Background(), reader.ID)
//	})
//
//	// Assert
//	assert.Nil(t, err)
//}
//
//func TestLibCardService_Create_ErrorLibCardAlreadyExists_Classic(t *testing.T) {
//	ctx := context.Background()
//	var err error
//
//	// Arrange
//	container, err := getContainerForClassicUnitTests()
//	if err != nil {
//		t.Fatal(err)
//	}
//	db, err := applyMigrations(container)
//	if err != nil {
//		t.Fatal(err)
//	}
//	reader := ommodels.NewReaderModelObjectMother().DefaultReader()
//	libCard := tdbmodels.NewLibCardModelBuilder().WithReaderID(reader.ID).Build()
//
//	_, _ = db.ExecContext(
//		ctx, `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
//		reader.ID,
//		reader.Fio,
//		reader.PhoneNumber,
//		reader.Age,
//		reader.Password,
//		reader.Role,
//	)
//
//	_, _ = db.ExecContext(
//		ctx, `insert into bs.lib_card values ($1, $2, $3, $4, $5, $6)`,
//		libCard.ID,
//		libCard.ReaderID,
//		libCard.LibCardNum,
//		libCard.Validity,
//		libCard.IssueDate,
//		libCard.ActionStatus,
//	)
//	libCardRepo := implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())
//	libCardService := impl.NewLibCardService(libCardRepo, logging.GetLoggerForTests())
//	defer func(db *sqlx.DB) {
//		if err = db.Close(); err != nil {
//			t.Fatalf("failed to close database connection: %v\n", err)
//		}
//	}(db)
//
//	defer func() {
//		if err = container.Terminate(ctx); err != nil {
//			t.Fatalf("failed to terminate container: %v\n", err)
//		}
//	}()
//
//	runner.Run(t, "error libCard already exist", func(t provider.T) {
//
//		// Act
//		err = libCardService.Create(context.Background(), reader.ID)
//	})
//
//	// Assert
//	assert.NotNil(t, err)
//	assert.Equal(t, errs.ErrLibCardAlreadyExist, err)
//}
