package serviceTests_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	mockrepo "github.com/nikitalystsev/BookSmart/internal/tests/unitTests/serviceTests/mocks"
	"github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/serviceTests/objectMother/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/serviceTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"testing"
)

type LibCardServiceTestsSuite struct {
	suite.Suite
}

func (lcsts *LibCardServiceTestsSuite) BeforeEach(t provider.T) {
	t.Tags("BookSmart", "libCard service", "suite", "steps")
}

/*
	Лондонский вариант
*/

func (lcsts *LibCardServiceTestsSuite) Test_Create_Success(t provider.T) {
	var (
		libCardService intf.ILibCardService
		readerID       uuid.UUID
		err            error
	)

	t.Title("Test Create LibCard Success")
	t.Description("The new libCard was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		libCardService = impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
		readerID = uuid.New()
		mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errs.ErrLibCardDoesNotExists)
		mockLibCardRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardService.Create(context.Background(), readerID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (lcsts *LibCardServiceTestsSuite) Test_Create_ErrorLibCardAlreadyExists(t provider.T) {
	var (
		libCardService intf.ILibCardService
		readerID       uuid.UUID
		err            error
	)

	t.Title("Test Create LibCard Error: libCard already exists")
	t.Description("A new library card was not created because the reader already has one")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		libCardService = impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
		readerID = uuid.New()
		libCard := tdbmodels.NewLibCardModelBuilder().WithReaderID(readerID).Build()
		mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(libCard, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardService.Create(context.Background(), readerID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrLibCardAlreadyExist, err)
	})
}

func (lcsts *LibCardServiceTestsSuite) Test_Update_Success(t provider.T) {
	var (
		libCardService intf.ILibCardService
		libCard        *models.LibCardModel
		err            error
	)

	t.Title("Test Update LibCard Success")
	t.Description("A new library card was successfully updated")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		libCardService = impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
		libCard = ommodels.NewLibCardModelObjectMother().ExpiredLibCard()
		mockLibCardRepo.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(libCard, nil)
		mockLibCardRepo.EXPECT().Update(gomock.Any(), libCard).Return(nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardService.Update(context.Background(), libCard)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (lcsts *LibCardServiceTestsSuite) Test_Update_ErrorCheckLibCardExistence(t provider.T) {
	var (
		libCardService intf.ILibCardService
		libCard        *models.LibCardModel
		err            error
	)

	t.Title("Test Update LibCard Error: check libCard existence")
	t.Description("The library card was not updated due to an error checking its existence")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		libCardService = impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
		libCard = ommodels.NewLibCardModelObjectMother().ExpiredLibCard()
		mockLibCardRepo.EXPECT().GetByNum(gomock.Any(), libCard.LibCardNum).Return(nil, errors.New("database error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardService.Update(context.Background(), libCard)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("database error"), err)
	})
}

func (lcsts *LibCardServiceTestsSuite) Test_GetByReaderID_Success(t provider.T) {
	var (
		libCardService intf.ILibCardService
		libCard        *models.LibCardModel
		findLibCard    *models.LibCardModel
		readerID       uuid.UUID
		err            error
	)

	t.Title("Test Get LibCard By Reader ID Success")
	t.Description("The library card was successfully retrieved by reader ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		libCardService = impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
		readerID = uuid.New()
		libCard = tdbmodels.NewLibCardModelBuilder().WithReaderID(readerID).Build()
		mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(libCard, nil)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findLibCard, err = libCardService.GetByReaderID(context.Background(), readerID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(libCard, findLibCard)
	})
}

func (lcsts *LibCardServiceTestsSuite) Test_GetByReaderID_ErrorLibCardDoesNotExists(t provider.T) {
	var (
		libCardService intf.ILibCardService
		findLibCard    *models.LibCardModel
		readerID       uuid.UUID
		err            error
	)

	t.Title("Test Get LibCard By Reader ID Error: libCard does not exists")
	t.Description("The library card was not successfully obtained because it does not exist")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		ctrl := gomock.NewController(t)
		mockLibCardRepo := mockrepo.NewMockILibCardRepo(ctrl)
		libCardService = impl.NewLibCardService(mockLibCardRepo, logging.GetLoggerForTests())
		readerID = uuid.New()
		mockLibCardRepo.EXPECT().GetByReaderID(gomock.Any(), readerID).Return(nil, errs.ErrLibCardDoesNotExists)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findLibCard, err = libCardService.GetByReaderID(context.Background(), readerID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrLibCardDoesNotExists, err)
		t.Assert().Nil(findLibCard)
	})
}

/*
	Классический вариант
*/

func (lcsts *LibCardServiceTestsSuite) Test_Create_Success_Classic(t provider.T) {
	var (
		container      *postgres.PostgresContainer
		db             *sqlx.DB
		libCardService intf.ILibCardService
		reader         *models.ReaderModel
		err            error
	)

	t.Title("Test Create LibCard Success")
	t.Description("The new libCard was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		container, err = unitTests.GetPostgresForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		db, err = unitTests.ApplyMigrations(container)
		if err != nil {
			t.Fatal(err)
		}
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		_, _ = db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)
		libCardRepo := implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())
		libCardService = impl.NewLibCardService(libCardRepo, logging.GetLoggerForTests())
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardService.Create(context.Background(), reader.ID)
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
}

func (lcsts *LibCardServiceTestsSuite) Test_Create_ErrorLibCardAlreadyExists_Classic(t provider.T) {
	var (
		container      *postgres.PostgresContainer
		db             *sqlx.DB
		libCardService intf.ILibCardService
		reader         *models.ReaderModel
		err            error
	)

	t.Title("Test Create LibCard Error: libCard already exists")
	t.Description("A new library card was not created because the reader already has one")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		container, err = unitTests.GetPostgresForClassicUnitTests()
		if err != nil {
			t.Fatal(err)
		}
		db, err = unitTests.ApplyMigrations(container)
		if err != nil {
			t.Fatal(err)
		}
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
		libCard := tdbmodels.NewLibCardModelBuilder().WithReaderID(reader.ID).Build()

		_, _ = db.ExecContext(
			context.Background(), `insert into bs.reader values ($1, $2, $3, $4, $5, $6)`,
			reader.ID,
			reader.Fio,
			reader.PhoneNumber,
			reader.Age,
			reader.Password,
			reader.Role,
		)

		_, _ = db.ExecContext(
			context.Background(), `insert into bs.lib_card values ($1, $2, $3, $4, $5, $6)`,
			libCard.ID,
			libCard.ReaderID,
			libCard.LibCardNum,
			libCard.Validity,
			libCard.IssueDate,
			libCard.ActionStatus,
		)
		libCardRepo := implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())
		libCardService = impl.NewLibCardService(libCardRepo, logging.GetLoggerForTests())
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardService.Create(context.Background(), reader.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errs.ErrLibCardAlreadyExist, err)
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
}

func TestLibCardServiceTestsSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(LibCardServiceTestsSuite))
}
