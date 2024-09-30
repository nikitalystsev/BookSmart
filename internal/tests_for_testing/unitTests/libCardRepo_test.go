package unitTests

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	implRepo "github.com/nikitalystsev/BookSmart-repo-postgres/impl"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/intfRepo"
	ommodels "github.com/nikitalystsev/BookSmart/internal/tests_for_testing/unitTests/objectMother/models"
	"github.com/nikitalystsev/BookSmart/pkg/logging"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

type LibCardRepoTestsSuite struct {
	suite.Suite
}

func (lcrts *LibCardRepoTestsSuite) BeforeEach(t provider.T) {
	t.Tags("BookSmart", "libCard repo", "suite", "steps")
}

func (lcrts *LibCardRepoTestsSuite) Test_Create_Success(t provider.T) {
	var (
		libCardRepo intfRepo.ILibCardRepo
		libCard     *models.LibCardModel
		mock        sqlxmock.Sqlmock
		db          *sqlx.DB
		err         error
	)

	t.Title("Test Create LibCard Success")
	t.Description("The new libCard was successfully created")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		libCard = ommodels.NewLibCardModelObjectMother().DefaultLibCard()
		libCardRepo = implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`insert into bs.lib_card values`).
			WithArgs(libCard.ID, libCard.ReaderID, libCard.LibCardNum,
				libCard.Validity, libCard.IssueDate, libCard.ActionStatus).
			WillReturnResult(sqlxmock.NewResult(1, 1))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardRepo.Create(context.Background(), libCard)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (lcrts *LibCardRepoTestsSuite) Test_Create_ErrorExecutionQuery(t provider.T) {
	var (
		libCardRepo intfRepo.ILibCardRepo
		libCard     *models.LibCardModel
		mock        sqlxmock.Sqlmock
		db          *sqlx.DB
		err         error
	)

	t.Title("Test Create LibCard Error: execution query error")
	t.Description("Execution query error")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}
		libCard = ommodels.NewLibCardModelObjectMother().DefaultLibCard()
		libCardRepo = implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`insert into bs.lib_card values`).
			WithArgs(libCard.ID, libCard.ReaderID, libCard.LibCardNum,
				libCard.Validity, libCard.IssueDate, libCard.ActionStatus).
			WillReturnError(errors.New("insert error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardRepo.Create(context.Background(), libCard)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("insert error"), err)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (lcrts *LibCardRepoTestsSuite) Test_GetByReaderID_Success(t provider.T) {
	var (
		libCardRepo intfRepo.ILibCardRepo
		libCard     *models.LibCardModel
		findLibCard *models.LibCardModel
		mock        sqlxmock.Sqlmock
		db          *sqlx.DB
		err         error
	)

	t.Title("Test Get LibCard By Reader ID Success")
	t.Description("The libCard was successfully retrieved by reader ID")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		libCard = ommodels.NewLibCardModelObjectMother().DefaultLibCard()
		libCardRepo = implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "reader_id", "lib_card_num", "validity", "issue_date", "action_status"}).
			AddRow(libCard.ID, libCard.ReaderID, libCard.LibCardNum, libCard.Validity, libCard.IssueDate, libCard.ActionStatus)

		mock.ExpectQuery(`select (.+) from bs.lib_card_view where (.+)`).
			WithArgs(libCard.ReaderID).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findLibCard, err = libCardRepo.GetByReaderID(context.Background(), libCard.ReaderID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(libCard, findLibCard)
	})
}

func (lcrts *LibCardRepoTestsSuite) Test_GetByReaderID_ErrorExecutionQuery(t provider.T) {
	var (
		libCardRepo intfRepo.ILibCardRepo
		libCard     *models.LibCardModel
		findLibCard *models.LibCardModel
		mock        sqlxmock.Sqlmock
		db          *sqlx.DB
		err         error
	)

	t.Title("Test Get LibCard By Reader ID Error: execution query error")
	t.Description("Execution query error")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}
		libCard = ommodels.NewLibCardModelObjectMother().DefaultLibCard()
		libCardRepo = implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.lib_card_view where (.+)`).
			WithArgs(libCard.ReaderID).WillReturnError(errors.New("query error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findLibCard, err = libCardRepo.GetByReaderID(context.Background(), libCard.ReaderID)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("query error"), err)
		t.Assert().Nil(findLibCard)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (lcrts *LibCardRepoTestsSuite) Test_GetByNum_Success(t provider.T) {
	var (
		libCardRepo intfRepo.ILibCardRepo
		libCard     *models.LibCardModel
		findLibCard *models.LibCardModel
		mock        sqlxmock.Sqlmock
		db          *sqlx.DB
		err         error
	)

	t.Title("Test Get LibCard By Num Success")
	t.Description("The libCard was successfully retrieved by num")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		libCard = ommodels.NewLibCardModelObjectMother().DefaultLibCard()
		libCardRepo = implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())

		rows := sqlxmock.NewRows([]string{"id", "reader_id", "lib_card_num", "validity", "issue_date", "action_status"}).
			AddRow(libCard.ID, libCard.ReaderID, libCard.LibCardNum, libCard.Validity, libCard.IssueDate, libCard.ActionStatus)

		mock.ExpectQuery(`select (.+) from bs.lib_card_view where (.+)`).
			WithArgs(libCard.LibCardNum).WillReturnRows(rows)
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findLibCard, err = libCardRepo.GetByNum(context.Background(), libCard.LibCardNum)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
		t.Assert().Equal(libCard, findLibCard)
	})
}

func (lcrts *LibCardRepoTestsSuite) Test_GetByNum_ErrorExecutionQuery(t provider.T) {
	var (
		libCardRepo intfRepo.ILibCardRepo
		libCard     *models.LibCardModel
		findLibCard *models.LibCardModel
		mock        sqlxmock.Sqlmock
		db          *sqlx.DB
		err         error
	)

	t.Title("Test Get LibCard By Num Error: execution query error")
	t.Description("Execution query error")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}
		libCard = ommodels.NewLibCardModelObjectMother().DefaultLibCard()
		libCardRepo = implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())

		mock.ExpectQuery(`select (.+) from bs.lib_card_view where (.+)`).
			WithArgs(libCard.LibCardNum).WillReturnError(errors.New("query error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		findLibCard, err = libCardRepo.GetByNum(context.Background(), libCard.LibCardNum)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("query error"), err)
		t.Assert().Nil(findLibCard)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func (lcrts *LibCardRepoTestsSuite) Test_Update_Success(t provider.T) {
	var (
		libCardRepo intfRepo.ILibCardRepo
		libCard     *models.LibCardModel
		mock        sqlxmock.Sqlmock
		db          *sqlx.DB
		err         error
	)

	t.Title("Test Update LibCard Success")
	t.Description("The libCard was successfully updated")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}

		libCard = ommodels.NewLibCardModelObjectMother().DefaultLibCard()
		libCardRepo = implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`update bs.lib_card set (.+) where (.+)`).
			WithArgs(libCard.ReaderID, libCard.LibCardNum, libCard.Validity,
				libCard.IssueDate, libCard.ActionStatus, libCard.ID).
			WillReturnResult(sqlxmock.NewResult(1, 1))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardRepo.Update(context.Background(), libCard)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().Nil(err)
	})
}

func (lcrts *LibCardRepoTestsSuite) Test_Update_ErrorExecutionQuery(t provider.T) {
	var (
		libCardRepo intfRepo.ILibCardRepo
		libCard     *models.LibCardModel
		mock        sqlxmock.Sqlmock
		db          *sqlx.DB
		err         error
	)

	t.Title("Test Update LibCard Error: execution query error")
	t.Description("Execution query error")

	t.WithNewStep("Arrange", func(sCtx provider.StepCtx) {
		if db, mock, err = sqlxmock.Newx(); err != nil {
			t.Fatal(err)
		}
		libCard = ommodels.NewLibCardModelObjectMother().DefaultLibCard()
		libCardRepo = implRepo.NewLibCardRepo(db, logging.GetLoggerForTests())

		mock.ExpectExec(`update bs.lib_card set (.+) where (.+)`).
			WithArgs(libCard.ReaderID, libCard.LibCardNum, libCard.Validity,
				libCard.IssueDate, libCard.ActionStatus, libCard.ID).
			WillReturnError(errors.New("update error"))
	})

	t.WithNewStep("Act", func(sCtx provider.StepCtx) {
		err = libCardRepo.Update(context.Background(), libCard)
	})

	t.WithNewStep("Assert", func(sCtx provider.StepCtx) {
		t.Assert().NotNil(err)
		t.Assert().Equal(errors.New("update error"), err)
		err = mock.ExpectationsWereMet()
		t.Assert().NoError(err)
	})
}

func TestLibCardRepoTestsSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(LibCardRepoTestsSuite))
}
