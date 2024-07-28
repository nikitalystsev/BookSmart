package postgres

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type ReservationRepo struct {
	db *sqlx.DB
}

func NewReservationRepo(db *sqlx.DB) intfRepo.IReservationRepo {
	return &ReservationRepo{db: db}
}

func (rr ReservationRepo) Create(ctx context.Context, reservation *models.ReservationModel) error {
	query := `INSERT INTO reservation VALUES ($1, $2, $3, $4, $5, $6)`

	if reservation.ID == uuid.Nil {
		reservation.ID = uuid.New()
	}

	_, err := rr.db.ExecContext(ctx, query, reservation.ID, reservation.ReaderID, reservation.BookID,
		reservation.IssueDate, reservation.ReturnDate, reservation.State)
	if err != nil {
		return fmt.Errorf("error creating reservation: %w", err)
	}

	return nil
}

func (rr ReservationRepo) GetByReaderAndBook(ctx context.Context, readerID, bookID uuid.UUID) (*models.ReservationModel, error) {
	query := `SELECT id, reader_id, book_id, issue_date, return_date, state FROM reservation WHERE reader_id = $1 AND book_id = $2`

	var reservation models.ReservationModel

	err := rr.db.GetContext(ctx, &reservation, query, readerID, bookID)
	if err != nil {
		return nil, fmt.Errorf("error getting reservation by reader and book: %w", err)
	}

	return &reservation, nil
}

func (rr ReservationRepo) GetByID(ctx context.Context, reservationID uuid.UUID) (*models.ReservationModel, error) {
	query := `SELECT id, reader_id, book_id, issue_date, return_date, state FROM reservation WHERE id = $1`

	var reservation models.ReservationModel

	err := rr.db.GetContext(ctx, &reservation, query, reservationID)
	if err != nil {
		return nil, fmt.Errorf("error getting reservation by id: %w", err)
	}

	return &reservation, nil
}

func (rr ReservationRepo) Update(ctx context.Context, reservation *models.ReservationModel) error {
	query := `UPDATE reservation SET issue_date = $1, return_date = $2, state = $3 WHERE id = $4`

	_, err := rr.db.ExecContext(ctx, query, reservation.IssueDate, reservation.ReturnDate, reservation.State, reservation.ID)
	if err != nil {
		return fmt.Errorf("error updating reservation: %w", err)
	}

	return nil
}

func (rr ReservationRepo) GetOverdueByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	query := `SELECT id, reader_id, book_id, issue_date, return_date, state FROM reservation WHERE reader_id = $1 AND return_date < $2`

	rows, err := rr.db.QueryxContext(ctx, query, readerID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error querying overdue reservations: %w", err)
	}
	defer func(rows *sqlx.Rows) {
		err = rows.Close()
		if err != nil {
			fmt.Printf("error")
		}
	}(rows)

	var reservations []*models.ReservationModel
	for rows.Next() {
		var reservation models.ReservationModel
		if err := rows.StructScan(&reservation); err != nil {
			return nil, fmt.Errorf("error scanning reservation: %w", err)
		}
		reservations = append(reservations, &reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}

	return reservations, nil
}

func (rr ReservationRepo) GetActiveByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	query := `SELECT id, reader_id, book_id, issue_date, return_date, state FROM reservation WHERE reader_id = $1 AND state != 'returned'`

	rows, err := rr.db.QueryxContext(ctx, query, readerID)
	if err != nil {
		return nil, fmt.Errorf("error querying active reservations: %w", err)
	}
	defer func(rows *sqlx.Rows) {
		err = rows.Close()
		if err != nil {
			fmt.Printf("error")
		}
	}(rows)

	var reservations []*models.ReservationModel
	for rows.Next() {
		var reservation models.ReservationModel
		if err := rows.StructScan(&reservation); err != nil {
			return nil, fmt.Errorf("error scanning reservation: %w", err)
		}
		reservations = append(reservations, &reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}

	return reservations, nil
}
