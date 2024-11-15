package integrationTests

import (
	"context"
	repomodels "github.com/nikitalystsev/BookSmart-repo-postgres/core/models"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	"github.com/nikitalystsev/BookSmart-services/pkg/auth"
	"github.com/nikitalystsev/BookSmart-services/pkg/hash"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (its *IntegrationTestSuite) TestReader_SignUp_Success(t provider.T) {
	t.Parallel()
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		findReaders   []*repomodels.ReaderModel
		err           error
	)

	t.Title("Integration Test SignUp Reader Success")
	t.Description("The new reader was successfully signUp")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		var tokenManager auth.ITokenManager
		if tokenManager, err = auth.NewTokenManager("signing_key"); err != nil {
			t.Fatal(err)
		}
		hasher := hash.NewPasswordHasher("salt_string")
		readerService = impl.NewReaderService(readerRepo, bookRepo, tokenManager, hasher,
			logging.GetLoggerForTests(), 10, 20)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("67890123456").Build()
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerService.SignUp(context.Background(), reader)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		query := `select * from bs.reader where id = $1`
		err = its.db.SelectContext(context.Background(), &findReaders, query, reader.ID)
		t.Assert().Nil(err)
		t.Assert().Len(findReaders, 1)
		t.Assert().Equal(reader.ID, findReaders[0].ID)
	})
}

func (its *IntegrationTestSuite) TestReader_SignIn_Success(t provider.T) {
	t.Parallel()
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		tokens        *models.Tokens
		err           error
	)

	t.Title("Integration Test SignIn Reader Success")
	t.Description("The new reader was successfully signIn")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		var tokenManager auth.ITokenManager
		if tokenManager, err = auth.NewTokenManager("signing_key"); err != nil {
			t.Fatal(err)
		}
		hasher := hash.NewPasswordHasher("salt_string")
		readerService = impl.NewReaderService(readerRepo, bookRepo, tokenManager, hasher,
			logging.GetLoggerForTests(), 10, 20)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("78901234567").Build()
		var password string
		if password, err = hasher.Hash(reader.Password); err != nil {
			t.Fatal(err)
		}
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			password,
			reader.Role,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		tokens, err = readerService.SignIn(context.Background(), reader.PhoneNumber, reader.Password)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().NotNil(tokens)
	})
}

func (its *IntegrationTestSuite) TestReader_GetByID_Success(t provider.T) {
	t.Parallel()
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		findReader    *models.ReaderModel
		err           error
	)

	t.Title("Integration Test Get Reader By ID Success")
	t.Description("The new reader was successfully getting by ID")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		var tokenManager auth.ITokenManager
		if tokenManager, err = auth.NewTokenManager("signing_key"); err != nil {
			t.Fatal(err)
		}
		hasher := hash.NewPasswordHasher("salt_string")
		readerService = impl.NewReaderService(readerRepo, bookRepo, tokenManager, hasher,
			logging.GetLoggerForTests(), 10, 20)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("89012345678").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReader, err = readerService.GetByID(context.Background(), reader.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(reader, findReader)
	})
}

func (its *IntegrationTestSuite) TestReader_GetByPhoneNumber_Success(t provider.T) {
	t.Parallel()
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		findReader    *models.ReaderModel
		err           error
	)

	t.Title("Integration Test Get Reader By ID Success")
	t.Description("The new reader was successfully getting by ID")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		var tokenManager auth.ITokenManager
		if tokenManager, err = auth.NewTokenManager("signing_key"); err != nil {
			t.Fatal(err)
		}
		hasher := hash.NewPasswordHasher("salt_string")
		readerService = impl.NewReaderService(readerRepo, bookRepo, tokenManager, hasher,
			logging.GetLoggerForTests(), 10, 20)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("90123456789").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findReader, err = readerService.GetByPhoneNumber(context.Background(), reader.PhoneNumber)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(reader, findReader)
	})
}

func (its *IntegrationTestSuite) TestReader_RefreshTokens_Success(t provider.T) {
	t.Parallel()
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		//findReader    *models.ReaderModel
		refreshToken string
		tokens       *models.Tokens
		err          error
	)

	t.Title("Integration Test Refresh Tokens Success")
	t.Description("The new reader was successfully refreshing tokens")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		var tokenManager auth.ITokenManager
		if tokenManager, err = auth.NewTokenManager("signing_key"); err != nil {
			t.Fatal(err)
		}
		hasher := hash.NewPasswordHasher("salt_string")
		readerService = impl.NewReaderService(readerRepo, bookRepo, tokenManager, hasher,
			logging.GetLoggerForTests(), 10, 20)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("01234567890").Build()
		var password string
		if password, err = hasher.Hash(reader.Password); err != nil {
			t.Fatal(err)
		}
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			password,
			reader.Role,
		)
		if err != nil {
			t.Fatal(err)
		}

		if refreshToken, err = tokenManager.NewRefreshToken(); err != nil {
			t.Fatal(err)
		}
		if err = its.client.Set(context.Background(), refreshToken, reader.ID.String(), 0).Err(); err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		tokens, err = readerService.RefreshTokens(context.Background(), refreshToken)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().NotNil(tokens)
	})
}

func (its *IntegrationTestSuite) TestReader_AddToFavorites_Success(t provider.T) {
	t.Parallel()
	var (
		readerService intf.IReaderService
		reader        *models.ReaderModel
		book          *models.BookModel
		err           error
	)

	t.Title("Integration Test Add To Favorites Success")
	t.Description("Book was successfully added to favorites")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		bookRepo := implRepo.NewBookRepo(its.db, logging.GetLoggerForTests())
		readerRepo := implRepo.NewReaderRepo(its.db, its.client, logging.GetLoggerForTests())
		var tokenManager auth.ITokenManager
		if tokenManager, err = auth.NewTokenManager("signing_key"); err != nil {
			t.Fatal(err)
		}
		hasher := hash.NewPasswordHasher("salt_string")
		readerService = impl.NewReaderService(readerRepo, bookRepo, tokenManager, hasher,
			logging.GetLoggerForTests(), 10, 20)
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("98765432109").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		if err != nil {
			t.Fatal(err)
		}
		book = ommodels.NewBookModelObjectMother().DefaultBook()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.book values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			book.ID,
			book.Title,
			book.Author,
			book.Publisher,
			book.CopiesNumber,
			book.Rarity,
			book.Genre,
			book.PublishingYear,
			book.Language,
			book.AgeLimit,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = readerService.AddToFavorites(context.Background(), reader.ID, book.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}
