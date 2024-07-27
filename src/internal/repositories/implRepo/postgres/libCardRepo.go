package postgres

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LibCardRepo struct {
	db *sqlx.DB
}

func NewLibCardRepo(db *sqlx.DB) intfRepo.ILibCardRepo {
	return &LibCardRepo{db: db}
}

func (lcr *LibCardRepo) Create(ctx context.Context, libCard *models.LibCardModel) error {
	query := `INSERT INTO lib_card VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := lcr.db.ExecContext(ctx, query, libCard.ID, libCard.ReaderID, libCard.LibCardNum,
		libCard.Validity, libCard.IssueDate, libCard.ActionStatus)

	if err != nil {
		return fmt.Errorf("error inserting libCard: %w", err)
	}

	return err
}

func (lcr *LibCardRepo) GetByReaderID(ctx context.Context, readerID uuid.UUID) (*models.LibCardModel, error) {
	query := `SELECT id, reader_id, lib_card_num, validity, issue_date, action_status FROM lib_card WHERE reader_id = $1`

	var libCard models.LibCardModel

	err := lcr.db.GetContext(ctx, &libCard, query, readerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("error retrieving lib card by reader ID: %w", err)
	}

	return &libCard, nil
}

func (lcr *LibCardRepo) GetByNum(ctx context.Context, libCardNum string) (*models.LibCardModel, error) {
	query := `SELECT id, reader_id, lib_card_num, validity, issue_date, action_status FROM lib_card WHERE lib_card_num = $1`

	var libCard models.LibCardModel

	err := lcr.db.GetContext(ctx, &libCard, query, libCardNum)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("error retrieving lib card by num: %w", err)
	}

	return &libCard, nil
}

func (lcr *LibCardRepo) Update(ctx context.Context, libCard *models.LibCardModel) error {
	query := `UPDATE lib_card SET issue_date = $1 WHERE id = $2
	`

	_, err := lcr.db.ExecContext(ctx, query, libCard.IssueDate, libCard.ID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error updating lib card: %v", err)
	}

	return nil
}
