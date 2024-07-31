package postgres

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errsRepo"
	"BookSmart/internal/repositories/intfRepo"
	"BookSmart/internal/services/implServices"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

type ReservationRepo struct {
	db     *sqlx.DB
	logger *logrus.Entry
}

func NewReservationRepo(db *sqlx.DB, logger *logrus.Entry) intfRepo.IReservationRepo {
	return &ReservationRepo{db: db, logger: logger}
}

func (rr *ReservationRepo) Create(ctx context.Context, reservation *models.ReservationModel) error {
	rr.logger.Infof("inserting reservation with ID: %s", reservation.ID)

	query := `INSERT INTO reservation VALUES ($1, $2, $3, $4, $5, $6)`

	rr.logger.Infof("executing query: %s", query)

	_, err := rr.db.ExecContext(ctx, query, reservation.ID, reservation.ReaderID, reservation.BookID,
		reservation.IssueDate, reservation.ReturnDate, reservation.State)
	if err != nil {
		rr.logger.Errorf("error inserting reservation: %v", err)
		return err
	}

	rr.logger.Infof("reservation with ID: %s inserted successfully", reservation.ID)

	return nil
}

func (rr *ReservationRepo) GetByReaderAndBook(ctx context.Context, readerID, bookID uuid.UUID) (*models.ReservationModel, error) {
	rr.logger.Infof("selecting reservation")

	query := `SELECT id, reader_id, book_id, issue_date, return_date, state FROM reservation WHERE reader_id = $1 AND book_id = $2`

	rr.logger.Infof("executing query: %s", query)

	var reservation models.ReservationModel
	err := rr.db.GetContext(ctx, &reservation, query, readerID, bookID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selected reservation")
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("no reservation found")
		return nil, errsRepo.ErrNotFound
	}

	if err = rr.checkExpiredReservation(ctx, &reservation); err != nil {
		return nil, err
	}

	rr.logger.Infof("successfully selected reservation: %v", reservation)

	return &reservation, nil
}

func (rr *ReservationRepo) GetByID(ctx context.Context, reservationID uuid.UUID) (*models.ReservationModel, error) {
	rr.logger.Infof("select reservation with ID: %s", reservationID)

	query := `SELECT id, reader_id, book_id, issue_date, return_date, state FROM reservation WHERE id = $1`

	rr.logger.Infof("executing query: %s", query)

	var reservation models.ReservationModel
	err := rr.db.GetContext(ctx, &reservation, query, reservationID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selected reservation by ID: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("no reservation found")
		return nil, errsRepo.ErrNotFound
	}

	if err = rr.checkExpiredReservation(ctx, &reservation); err != nil {
		return nil, err
	}

	rr.logger.Infof("successfully selected reservation: %v", reservation)

	return &reservation, nil
}

func (rr *ReservationRepo) Update(ctx context.Context, reservation *models.ReservationModel) error {
	rr.logger.Infof("updating reservation with readerID: %s", reservation.ID)

	query := `UPDATE reservation SET issue_date = $1, return_date = $2, state = $3 WHERE id = $4`

	rr.logger.Infof("executing query: %s", query)

	_, err := rr.db.ExecContext(ctx, query, reservation.IssueDate, reservation.ReturnDate, reservation.State, reservation.ID)
	if err != nil {
		rr.logger.Errorf("error updating reservation")
		return err
	}

	rr.logger.Infof("successfully updated reservation with ID: %s", reservation.ID)

	return nil
}

func (rr *ReservationRepo) GetExpiredByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	rr.logger.Infof("getting expired reservation with readerID: %s", readerID)

	query := `SELECT id, reader_id, book_id, issue_date, return_date, state FROM reservation WHERE reader_id = $1 AND return_date < $2`

	rr.logger.Infof("executing query: %s", query)

	rows, err := rr.db.QueryxContext(ctx, query, readerID, time.Now())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting expired books: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("no found expired books")
		return nil, errsRepo.ErrNotFound
	}

	defer func(rows *sqlx.Rows) {
		err = rows.Close()
		if err != nil {
			rr.logger.Errorf("error closing rows: %v", err)
			fmt.Printf("error closing rows: %v", err)
		}
	}(rows)

	var reservations []*models.ReservationModel
	for rows.Next() {
		var reservation models.ReservationModel
		if err = rows.StructScan(&reservation); err != nil {
			rr.logger.Errorf("error scanning reservations row: %v", err)
			return nil, err
		}
		if err = rr.checkExpiredReservation(ctx, &reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}

	if err = rows.Err(); err != nil {
		rr.logger.Errorf("rows iteration error: %v", err)
		return nil, err
	}

	rr.logger.Infof("successfully found %d expired reservations", len(reservations))

	return reservations, nil
}

func (rr *ReservationRepo) GetActiveByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	rr.logger.Infof("getting active reservation with readerID: %s", readerID)

	query := `SELECT id, reader_id, book_id, issue_date, return_date, state FROM reservation WHERE reader_id = $1 AND state != 'Closed'`

	rr.logger.Infof("executing query: %s", query)

	rows, err := rr.db.QueryxContext(ctx, query, readerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting active books: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("no found active books")
		return nil, errsRepo.ErrNotFound
	}

	defer func(rows *sqlx.Rows) {
		err = rows.Close()
		if err != nil {
			rr.logger.Errorf("error closing rows: %v", err)
			fmt.Printf("error closing rows: %v", err)
		}
	}(rows)

	var reservations []*models.ReservationModel
	for rows.Next() {
		var reservation models.ReservationModel
		if err = rows.StructScan(&reservation); err != nil {
			rr.logger.Errorf("error scanning reservations row: %v", err)
			return nil, err
		}
		if err = rr.checkExpiredReservation(ctx, &reservation); err != nil {
			return nil, err
		}
		reservations = append(reservations, &reservation)
	}

	if err = rows.Err(); err != nil {
		rr.logger.Errorf("rows iteration error: %v", err)
		return nil, err
	}

	rr.logger.Infof("successfully found %d active reservations", len(reservations))

	return reservations, nil
}

func (rr *ReservationRepo) checkExpiredReservation(ctx context.Context, reservation *models.ReservationModel) error {
	if !rr.isExpiredReservation(reservation) {
		return nil
	}

	if reservation.State == implServices.ReservationExpired {
		rr.logger.Infof("reservation is already expired")
		return nil
	}

	reservation.State = implServices.ReservationExpired
	if err := rr.updateReservationStatus(ctx, reservation); err != nil {
		return err
	}

	return nil
}

func (rr *ReservationRepo) isExpiredReservation(reservation *models.ReservationModel) bool {
	rr.logger.Infof("check reservation status: %s", reservation.ID)

	if reservation.State == implServices.ReservationClosed || !time.Now().After(reservation.ReturnDate) {
		rr.logger.Infof("reservation does not expired")
		return false
	}

	rr.logger.Infof("reservation is expired")

	return true
}

func (rr *ReservationRepo) updateReservationStatus(ctx context.Context, reservation *models.ReservationModel) error {
	rr.logger.Infof("updating reservation status: %s", reservation.ID)

	query := `UPDATE reservation SET state = $1 WHERE id = $2`

	rr.logger.Infof("executing query: %s", query)

	_, err := rr.db.ExecContext(ctx, query, reservation.State, reservation.ID)
	if err != nil {
		rr.logger.Printf("error updating reservation status: %v", err)
		return err
	}

	rr.logger.Infof("successfully updated reservation status")

	return nil
}
