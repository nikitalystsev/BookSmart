package integrationTests

import (
	"context"
	repomodels "github.com/nikitalystsev/BookSmart-repo-postgres/core/models"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/impl"
	"github.com/nikitalystsev/BookSmart-services/intf"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	tdbmodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/testDataBuilder/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"time"
)

func (its *IntegrationTestSuite) TestLibCard_Create_Success(t provider.T) {
	t.Parallel()
	var (
		libCardService intf.ILibCardService
		reader         *models.ReaderModel
		findLibCards   []*repomodels.LibCardModel
		err            error
	)

	t.Title("Integration Test Create LibCard Success")
	t.Description("The new libCard was successfully created")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		libCardRepo := implRepo.NewLibCardRepo(its.db, logging.GetLoggerForTests())
		libCardService = impl.NewLibCardService(libCardRepo, logging.GetLoggerForTests())
		reader = ommodels.NewReaderModelObjectMother().DefaultReader()
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
		err = libCardService.Create(context.Background(), reader.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		query := `select * from bs.lib_card where reader_id = $1`
		err = its.db.SelectContext(context.Background(), &findLibCards, query, reader.ID)
		t.Assert().Nil(err)
		t.Assert().Len(findLibCards, 1)
		t.Assert().Equal(reader.ID, findLibCards[0].ReaderID)
	})
}

func (its *IntegrationTestSuite) TestLibCard_Update_Success(t provider.T) {
	t.Parallel()
	var (
		libCardService intf.ILibCardService
		reader         *models.ReaderModel
		libCard        *models.LibCardModel
		findLibCards   []*repomodels.LibCardModel
		err            error
	)

	t.Title("Integration Test Update LibCard Success")
	t.Description("The new libCard was successfully updated")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		libCardRepo := implRepo.NewLibCardRepo(its.db, logging.GetLoggerForTests())
		libCardService = impl.NewLibCardService(libCardRepo, logging.GetLoggerForTests())
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("12345678901").Build()
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
		libCard = tdbmodels.NewLibCardModelBuilder().
			WithReaderID(reader.ID).
			WithIssueDate(time.Now().AddDate(0, 0, -370)).
			WithActionStatus(false).Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.lib_card values ($1, $2, $3, $4, $5, $6)`,
			libCard.ID,
			libCard.ReaderID,
			libCard.LibCardNum,
			libCard.Validity,
			libCard.IssueDate,
			libCard.ActionStatus,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardService.Update(context.Background(), libCard)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		query := `select * from bs.lib_card where reader_id = $1`
		err = its.db.SelectContext(context.Background(), &findLibCards, query, reader.ID)
		t.Assert().Nil(err)
		t.Assert().Len(findLibCards, 1)
		t.Assert().Equal(reader.ID, findLibCards[0].ReaderID)
		t.Assert().Equal(true, findLibCards[0].ActionStatus)
	})
}

func (its *IntegrationTestSuite) TestLibCard_GetByReaderID_Success(t provider.T) {
	t.Parallel()
	var (
		libCardService intf.ILibCardService
		reader         *models.ReaderModel
		libCard        *models.LibCardModel
		findLibCard    *models.LibCardModel
		err            error
	)

	t.Title("Integration Test Get LibCard By Reader ID Success")
	t.Description("LibCard was successfully getting by reader ID")
	if isUnitTestsFailed() {
		t.Skip()
	}

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		libCardRepo := implRepo.NewLibCardRepo(its.db, logging.GetLoggerForTests())
		libCardService = impl.NewLibCardService(libCardRepo, logging.GetLoggerForTests())
		reader = tdbmodels.NewReaderModelBuilder().
			WithPhoneNumber("23456789012").Build()
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
		libCard = tdbmodels.NewLibCardModelBuilder().
			WithReaderID(reader.ID).
			WithLibCardNum("234567890123").Build()
		_, err = its.db.ExecContext(
			context.Background(), `insert into bs.lib_card values ($1, $2, $3, $4, $5, $6)`,
			libCard.ID,
			libCard.ReaderID,
			libCard.LibCardNum,
			libCard.Validity,
			libCard.IssueDate,
			libCard.ActionStatus,
		)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findLibCard, err = libCardService.GetByReaderID(context.Background(), reader.ID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(libCard.ID, findLibCard.ID)
		t.Assert().Equal(libCard.LibCardNum, findLibCard.LibCardNum)
		t.Assert().Equal(libCard.ReaderID, findLibCard.ReaderID)
	})
}
