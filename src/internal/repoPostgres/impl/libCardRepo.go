package impl

import (
	"BookSmart-services/core/models"
	"BookSmart-services/errs"
	"BookSmart-services/intfRepo"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type LibCardRepo struct {
	db     *sqlx.DB
	logger *logrus.Entry
}

func NewLibCardRepo(db *sqlx.DB, logger *logrus.Entry) intfRepo.ILibCardRepo {
	return &LibCardRepo{db: db, logger: logger}
}

func (lcr *LibCardRepo) Create(ctx context.Context, libCard *models.LibCardModel) error {
	lcr.logger.Infof("inserting libCard with ID: %s", libCard.ID)

	query := `insert into bs.lib_card values ($1, $2, $3, $4, $5, $6)`

	result, err := lcr.db.ExecContext(ctx, query, libCard.ID, libCard.ReaderID, libCard.LibCardNum,
		libCard.Validity, libCard.IssueDate, libCard.ActionStatus)

	if err != nil {
		lcr.logger.Errorf("error inserting libCard: %v", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		lcr.logger.Errorf("error inserting libCard: %v", err)
		return err
	}
	if rows != 1 {
		lcr.logger.Errorf("error inserting libCard: expected 1 row affected, got %d", rows)
		return errors.New("libCardRepo.Create: expected 1 row affected")
	}

	lcr.logger.Infof("inserted libCard with ID: %s", libCard.ID)

	return err
}

func (lcr *LibCardRepo) GetByReaderID(ctx context.Context, readerID uuid.UUID) (*models.LibCardModel, error) {
	lcr.logger.Infof("selecting libCard with readerID: %s", readerID)

	query := `select * from bs.lib_card_view where reader_id = $1`

	var libCard models.LibCardModel
	err := lcr.db.GetContext(ctx, &libCard, query, readerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		lcr.logger.Errorf("error selecting libCard: %v", err)
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		lcr.logger.Warnf("libCard with this readerID not found: %v", readerID)
		return nil, errs.ErrLibCardDoesNotExists
	}

	lcr.logger.Infof("selected libCard with readerID: %s", readerID)

	return &libCard, nil
}

func (lcr *LibCardRepo) GetByNum(ctx context.Context, libCardNum string) (*models.LibCardModel, error) {
	lcr.logger.Infof("selecting libCard with num: %s", libCardNum)

	query := `select * from bs.lib_card_view where lib_card_num = $1`

	lcr.logger.Infof("executing query: %s", query)

	var libCard models.LibCardModel
	err := lcr.db.GetContext(ctx, &libCard, query, libCardNum)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		lcr.logger.Errorf("error selected libCard with num: %v", err)
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		lcr.logger.Warnf("libCard with this num not found: %v", libCardNum)
		return nil, errs.ErrLibCardDoesNotExists
	}

	lcr.logger.Infof("selected libCard with num: %s", libCardNum)

	return &libCard, nil
}

func (lcr *LibCardRepo) Update(ctx context.Context, libCard *models.LibCardModel) error {
	lcr.logger.Infof("updating libCard with ID: %s", libCard.ID)

	query := `update bs.lib_card set issue_date = $1, action_status = $2 where id = $3`

	result, err := lcr.db.ExecContext(ctx, query, libCard.IssueDate, libCard.ActionStatus, libCard.ID)
	if err != nil {
		lcr.logger.Errorf("error updating libCard: %v", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		lcr.logger.Errorf("error updating libCard: %v", err)
		return err
	}
	if rows != 1 {
		lcr.logger.Errorf("error updating libCard: expected 1 row affected, got %d", rows)
		return errors.New("libCardRepo.Update: expected 1 row affected")
	}

	lcr.logger.Infof("updated libCard with ID: %s", libCard.ID)

	return nil
}
