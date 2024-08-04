package postgres

import (
	"BookSmart-repositories/errs"
	"BookSmart-repositories/intf"
	"BookSmart-services/models"
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

func NewLibCardRepo(db *sqlx.DB, logger *logrus.Entry) intf.ILibCardRepo {
	return &LibCardRepo{db: db, logger: logger}
}

func (lcr *LibCardRepo) Create(ctx context.Context, libCard *models.LibCardModel) error {
	lcr.logger.Infof("inserting libCard with ID: %s", libCard.ID)

	query := `INSERT INTO bs.lib_card VALUES ($1, $2, $3, $4, $5, $6)`

	lcr.logger.Infof("executing query: %s", query)

	_, err := lcr.db.ExecContext(ctx, query, libCard.ID, libCard.ReaderID, libCard.LibCardNum,
		libCard.Validity, libCard.IssueDate, libCard.ActionStatus)

	if err != nil {
		lcr.logger.Errorf("error inserting libCard: %v", err)
		return err
	}

	lcr.logger.Infof("libCard with ID: %s inserted successfully", libCard.ID)

	return err
}

func (lcr *LibCardRepo) GetByReaderID(ctx context.Context, readerID uuid.UUID) (*models.LibCardModel, error) {
	lcr.logger.Infof("select libCard with readerID: %s", readerID)

	query := `SELECT id, reader_id, lib_card_num, validity, issue_date, action_status FROM bs.lib_card WHERE reader_id = $1`

	lcr.logger.Infof("executing query: %s", query)

	var libCard models.LibCardModel
	err := lcr.db.GetContext(ctx, &libCard, query, readerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		lcr.logger.Errorf("error selected libCard by readerID: %v", readerID)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		lcr.logger.Warnf("no libCard found by readerID: %v", readerID)
		return nil, errs.ErrNotFound
	}

	lcr.logger.Infof("successfully selected libCard: %v", libCard)

	return &libCard, nil
}

func (lcr *LibCardRepo) GetByNum(ctx context.Context, libCardNum string) (*models.LibCardModel, error) {
	lcr.logger.Infof("select libCard with libCardNum: %s", libCardNum)

	query := `SELECT id, reader_id, lib_card_num, validity, issue_date, action_status FROM bs.lib_card WHERE lib_card_num = $1`

	lcr.logger.Infof("executing query: %s", query)

	var libCard models.LibCardModel
	err := lcr.db.GetContext(ctx, &libCard, query, libCardNum)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		lcr.logger.Errorf("error selected libCard by libCardNum: %v", libCardNum)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		lcr.logger.Warnf("no libCard found by libCardNum: %v", libCardNum)
		return nil, errs.ErrNotFound
	}

	lcr.logger.Infof("successfully selected libCard: %v", libCard)

	return &libCard, nil
}

func (lcr *LibCardRepo) Update(ctx context.Context, libCard *models.LibCardModel) error {
	lcr.logger.Infof("updating libCard with ID: %s", libCard.ID)

	query := `UPDATE bs.lib_card SET issue_date = $1 WHERE id = $2`

	lcr.logger.Infof("executing query: %s", query)

	_, err := lcr.db.ExecContext(ctx, query, libCard.IssueDate, libCard.ID)
	if err != nil {
		lcr.logger.Errorf("error updating libCard: %v", err)
		return err
	}

	lcr.logger.Infof("successfully updated book copies for ID: %s", libCard.ID)

	return nil
}
