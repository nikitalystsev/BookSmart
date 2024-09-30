package serviceTests

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	"github.com/nikitalystsev/BookSmart-services/pkg/auth"
	"github.com/nikitalystsev/BookSmart-services/pkg/hash"
	mockrepo "github.com/nikitalystsev/BookSmart/internal/tests/unitTests/serviceTests/mocks"
	omdto "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/serviceTests/objectMother/dto"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/serviceTests/objectMother/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	testredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"log"
	"testing"
)

type ReaderServiceTestsSuite struct {
	suite.Suite
}

func (rsts *ReaderServiceTestsSuite) BeforeEach(t provider.T) {
	t.Tags("BookSmart", "reader service", "suite", "steps")
}

/*
	Лондонский вариант
*/

func (rsts *ReaderServiceTestsSuite) Test_SignUp_Success(t provider.T) {
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		err           error
	)

	t.Title("Test Reader SingUp Success")
	t.Description("The new reader was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
		readerService = impl.NewReaderService(
			mockReaderRepo, nil, nil,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		mockReaderRepo.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
		mockHasher.EXPECT().Hash(reader.Password).Return("hashed password", nil)
		mockReaderRepo.EXPECT().Create(gomock.Any(), reader).Return(nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerService.SignUp(context.Background(), reader)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (rsts *ReaderServiceTestsSuite) Test_SingUp_ErrorCheckReaderExistence(t provider.T) {
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		err           error
	)

	t.Title("Test Reader SingUp Error: check reader existence")
	t.Description("The reader was not created successfully due to an error checking its existence")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
		readerService = impl.NewReaderService(
			mockReaderRepo, nil, nil,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		mockReaderRepo.EXPECT().GetByPhoneNumber(gomock.Any(), reader.PhoneNumber).Return(nil, errors.New("database error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerService.SignUp(context.Background(), reader)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("database error"), err)
	})
}

func (rsts *ReaderServiceTestsSuite) Test_SingIn_Success(t provider.T) {
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		readerDTO     *dto.SignInInputDTO
		tokens        *models.Tokens
		err           error
	)

	t.Title("Test Reader SingIn Success")
	t.Description("The reader was successfully sign in")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
		mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
		readerService = impl.NewReaderService(
			mockReaderRepo, nil, mockTokenManager,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		readerDTO = omdto.NewReaderSignInDTOObjectMother().DefaultReaderSignInDTO()
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		mockReaderRepo.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
		mockHasher.EXPECT().Compare(reader.Password, readerDTO.Password).Return(nil)
		mockTokenManager.EXPECT().NewJWT(reader.ID, reader.Role, gomock.Any()).Return("accessToken", nil)
		mockTokenManager.EXPECT().NewRefreshToken().Return("refreshToken", nil)
		mockReaderRepo.EXPECT().SaveRefreshToken(gomock.Any(), reader.ID, "refreshToken", gomock.Any()).Return(nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		tokens, err = readerService.SignIn(context.Background(), readerDTO.PhoneNumber, reader.Password)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().NotNil(tokens)
	})
}

func (rsts *ReaderServiceTestsSuite) Test_SingIn_ErrorReaderNotFound(t provider.T) {
	var (
		readerService intf.IReaderService
		readerDTO     *dto.SignInInputDTO
		tokens        *models.Tokens
		err           error
	)

	t.Title("Test Reader SingIn Error: reader not found")
	t.Description("The reader was not found when attempting to log in")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
		mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
		readerService = impl.NewReaderService(
			mockReaderRepo, nil, mockTokenManager,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		readerDTO = omdto.NewReaderSignInDTOObjectMother().DefaultReaderSignInDTO()
		mockReaderRepo.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		tokens, err = readerService.SignIn(context.Background(), readerDTO.PhoneNumber, readerDTO.Password)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrReaderDoesNotExists, err)
		t.Assert().Nil(tokens)
	})
}

func (rsts *ReaderServiceTestsSuite) Test_GetByPhoneNumber_Success(t provider.T) {
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		readerDTO     *dto.SignInInputDTO
		findReader    *models.ReaderModel
		err           error
	)

	t.Title("Test Get Reader By PhoneNumber Success")
	t.Description("The reader was successfully get by phoneNumber")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
		mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
		readerService = impl.NewReaderService(
			mockReaderRepo, nil, mockTokenManager,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		readerDTO = omdto.NewReaderSignInDTOObjectMother().DefaultReaderSignInDTO()
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		mockReaderRepo.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(reader, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReader, err = readerService.GetByPhoneNumber(context.Background(), readerDTO.PhoneNumber)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().NotNil(findReader)
		t.Assert().Equal(reader, findReader)
	})
}

func (rsts *ReaderServiceTestsSuite) Test_GetByPhoneNumber_ErrorReaderDoesNotExists(t provider.T) {
	var (
		readerService intf.IReaderService
		readerDTO     *dto.SignInInputDTO
		findReader    *models.ReaderModel
		err           error
	)

	t.Title("Test Get Reader By PhoneNumber Error: reader does not exists")
	t.Description("The reader was not found when trying to get information about him")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
		mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
		readerService = impl.NewReaderService(
			mockReaderRepo, nil, mockTokenManager,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		readerDTO = omdto.NewReaderSignInDTOObjectMother().DefaultReaderSignInDTO()
		mockReaderRepo.EXPECT().GetByPhoneNumber(gomock.Any(), readerDTO.PhoneNumber).Return(nil, errs.ErrReaderDoesNotExists)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {

		findReader, err = readerService.GetByPhoneNumber(context.Background(), readerDTO.PhoneNumber)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Nil(findReader)
		t.Assert().Equal(errs.ErrReaderDoesNotExists, err)
	})
}

func (rsts *ReaderServiceTestsSuite) Test_RefreshTokens_Success(t provider.T) {
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		tokens        *models.Tokens
		err           error
	)

	t.Title("Test Refresh Tokens Success")
	t.Description("Tokens have been successfully updated")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
		mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
		readerService = impl.NewReaderService(
			mockReaderRepo, nil, mockTokenManager,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		mockReaderRepo.EXPECT().GetByRefreshToken(gomock.Any(), "oldRefreshToken").Return(reader, nil)
		mockTokenManager.EXPECT().NewJWT(reader.ID, reader.Role, gomock.Any()).Return("newAccessToken", nil)
		mockTokenManager.EXPECT().NewRefreshToken().Return("newRefreshToken", nil)
		mockReaderRepo.EXPECT().SaveRefreshToken(gomock.Any(), reader.ID, "newRefreshToken", gomock.Any()).Return(nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		tokens, err = readerService.RefreshTokens(context.Background(), "oldRefreshToken")
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().NotNil(tokens)
	})
}

func (rsts *ReaderServiceTestsSuite) Test_RefreshTokens_ErrorFailedGettingReaderByRefreshToken(t provider.T) {
	var (
		readerService intf.IReaderService
		tokens        *models.Tokens
		err           error
	)

	t.Title("Test Refresh Tokens Error: failed to get reader by refresh token")
	t.Description("Tokens were successfully refreshed due to error getting reader on refresh token")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
		mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
		readerService = impl.NewReaderService(
			mockReaderRepo, nil, mockTokenManager,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		mockReaderRepo.EXPECT().GetByRefreshToken(gomock.Any(), "oldRefreshToken").Return(nil, errs.ErrReaderDoesNotExists)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		tokens, err = readerService.RefreshTokens(context.Background(), "oldRefreshToken")
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Nil(tokens)
		t.Assert().Equal(errs.ErrReaderDoesNotExists, err)
	})
}

func (rsts *ReaderServiceTestsSuite) Test_AddToFavorites_Success(t provider.T) {
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		book          *models.BookModel
		err           error
	)

	t.Title("Test Add Book To Favorites Success")
	t.Description("Book was successfully added to favorites")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
		mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
		readerService = impl.NewReaderService(
			mockReaderRepo, mockBookRepo, mockTokenManager,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		mockReaderRepo.EXPECT().GetByID(gomock.Any(), reader.ID).Return(reader, nil)
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
		mockReaderRepo.EXPECT().IsFavorite(gomock.Any(), reader.ID, book.ID).Return(false, nil)
		mockReaderRepo.EXPECT().AddToFavorites(gomock.Any(), reader.ID, book.ID).Return(nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerService.AddToFavorites(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (rsts *ReaderServiceTestsSuite) Test_AddToFavorites_ErrorBookAlreadyInFavorites(t provider.T) {
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		book          *models.BookModel
		err           error
	)

	t.Title("Test Add Book To Favorites Error: book already in favorites")
	t.Description("The book was not added to your favorites because it is already there")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockReaderRepo := mockrepo.NewMockIReaderRepo(ctrl)
		mockBookRepo := mockrepo.NewMockIBookRepo(ctrl)
		mockHasher := mockrepo.NewMockIPasswordHasher(ctrl)
		mockTokenManager := mockrepo.NewMockITokenManager(ctrl)
		readerService = impl.NewReaderService(
			mockReaderRepo, mockBookRepo, mockTokenManager,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		mockReaderRepo.EXPECT().GetByID(gomock.Any(), reader.ID).Return(reader, nil)
		mockBookRepo.EXPECT().GetByID(gomock.Any(), book.ID).Return(book, nil)
		mockReaderRepo.EXPECT().IsFavorite(gomock.Any(), reader.ID, book.ID).Return(true, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerService.AddToFavorites(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrBookAlreadyIsFavorite, err)
	})
}

/*
	Классический вариант
*/

func (rsts *ReaderServiceTestsSuite) Test_SignUp_Success_Classic(t provider.T) {
	var (
		container      *postgres.PostgresContainer
		redisContainer *testredis.RedisContainer
		client         *redis.Client
		db             *sqlx.DB
		readerService  intf.IReaderService
		reader         *models.ReaderModel
		err            error
	)

	t.Title("Test Reader SingUp Success Classic")
	t.Description("The new reader was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		container, err = getPostgresForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		db, err = applyMigrations(container)
		if err != nil {
			t.Fatal(err)
		}
		redisContainer, err = getRedisForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		client, err = getRedisClientForClassicUnitTests(redisContainer)
		if err != nil {
			t.Fatal(err)
		}
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		mockReaderRepo := implRepo.NewReaderRepo(db, client, logging.GetLoggerForTests())
		mockBookRepo := implRepo.NewBookRepo(db, logging.GetLoggerForTests())
		mockTokenManager, err := auth.NewTokenManager("signing_string")
		if err != nil {
			t.Fatal(err)
		}
		mockHasher := hash.NewPasswordHasher("salt_string")
		readerService = impl.NewReaderService(
			mockReaderRepo, mockBookRepo, mockTokenManager,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerService.SignUp(context.Background(), reader)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})

	defer func(db *sqlx.DB) {
		if err = db.Close(); err != nil {
			t.Fatalf("failed to close database connection: %v\n", err)
		}
	}(db)

	defer func() {
		if err = container.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %v\n", err)
		}
	}()
	defer func() {
		if err = redisContainer.Terminate(context.Background()); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
}

func (rsts *ReaderServiceTestsSuite) Test_SingUp_ErrorCheckReaderExistence_Classic(t provider.T) {
	var (
		container      *postgres.PostgresContainer
		redisContainer *testredis.RedisContainer
		client         *redis.Client
		db             *sqlx.DB
		readerService  intf.IReaderService
		reader         *models.ReaderModel
		err            error
	)

	t.Title("Test Reader SingUp Error: check reader existence")
	t.Description("The reader was not created successfully due to an error checking its existence")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		container, err = getPostgresForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		db, err = applyMigrations(container)
		if err != nil {
			t.Fatal(err)
		}
		redisContainer, err = getRedisForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		client, err = getRedisClientForClassicUnitTests(redisContainer)
		if err != nil {
			t.Fatal(err)
		}
		if err = db.Close(); err != nil {
			t.Fatalf("failed to close database connection: %v\n", err)
		}
		mockReaderRepo := implRepo.NewReaderRepo(db, client, logging.GetLoggerForTests())
		mockBookRepo := implRepo.NewBookRepo(db, logging.GetLoggerForTests())
		mockHasher := hash.NewPasswordHasher("salt_string")
		mockTokenManager, err := auth.NewTokenManager("signing_string")
		if err != nil {
			t.Fatal(err)
		}
		readerService = impl.NewReaderService(
			mockReaderRepo, mockBookRepo, mockTokenManager,
			mockHasher, logging.GetLoggerForTests(),
			1, 2,
		)
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerService.SignUp(context.Background(), reader)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Error(err)
	})

	defer func(db *sqlx.DB) {
		if err = db.Close(); err != nil {
			t.Fatalf("failed to close database connection: %v\n", err)
		}
	}(db)

	defer func() {
		if err = container.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %v\n", err)
		}
	}()

	defer func() {
		if err = redisContainer.Terminate(context.Background()); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
}

func TestReaderServiceTestsSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(ReaderServiceTestsSuite))
}
