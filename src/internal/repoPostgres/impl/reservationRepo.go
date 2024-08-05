package impl

import (
	"BookSmart-repositories/errs"
	"BookSmart-services/core/models"
	"BookSmart-services/impl"
	"BookSmart-services/intfRepo"
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

	query := `INSERT INTO bs.reservation VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := rr.db.ExecContext(ctx, query, reservation.ID, reservation.ReaderID, reservation.BookID,
		reservation.IssueDate, reservation.ReturnDate, reservation.State)
	if err != nil {
		rr.logger.Errorf("error inserting reservation: %v", err)
		return err
	}

	rr.logger.Infof("inserted reservation with ID: %s", reservation.ID)

	return nil
}

func (rr *ReservationRepo) GetByReaderAndBook(ctx context.Context, readerID, bookID uuid.UUID) (*models.ReservationModel, error) {
	rr.logger.Infof("selecting reservation with readerID и bookID: %s и %s", readerID, bookID)

	query := `SELECT * FROM bs.reservation_view WHERE reader_id = $1 AND book_id = $2`

	var reservation models.ReservationModel
	err := rr.db.GetContext(ctx, &reservation, query, readerID, bookID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting reservation: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("reservation with this readerID и bookID not found: %s и %s", readerID, bookID)
		return nil, errs.ErrNotFound
	}

	rr.logger.Infof("selected reservation with readerID и bookID: %s и %s", readerID, bookID)

	return &reservation, nil
}

func (rr *ReservationRepo) GetByID(ctx context.Context, ID uuid.UUID) (*models.ReservationModel, error) {
	rr.logger.Infof("selecting reservation with ID: %s", ID)

	query := `SELECT * FROM bs.reservation_view WHERE id = $1`

	var reservation models.ReservationModel
	err := rr.db.GetContext(ctx, &reservation, query, ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting reservation with ID: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("rewservation with this ID not found: %s", ID)
		return nil, errs.ErrNotFound
	}

	rr.logger.Infof("selected reservation with ID: %s", ID)

	return &reservation, nil
}

func (rr *ReservationRepo) Update(ctx context.Context, reservation *models.ReservationModel) error {
	rr.logger.Infof("updating reservation with ID: %s", reservation.ID)

	query := `UPDATE bs.reservation SET issue_date = $1, return_date = $2, state = $3 WHERE id = $4`

	_, err := rr.db.ExecContext(ctx, query, reservation.IssueDate, reservation.ReturnDate, reservation.State, reservation.ID)
	if err != nil {
		rr.logger.Errorf("error updating reservation with ID: %v", err)
		return err
	}

	rr.logger.Infof("updated reservation with ID: %s", reservation.ID)

	return nil
}

func (rr *ReservationRepo) GetExpiredByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	rr.logger.Infof("selecting expired reservations with readerID: %s", readerID)

	query := `SELECT * FROM bs.reservation_view WHERE reader_id = $1 AND return_date < $2`

	var reservations []*models.ReservationModel
	err := rr.db.SelectContext(ctx, &reservations, query, readerID, time.Now())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting expired reservations: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("expired reservations with this readerID not found: %s", readerID)
		return nil, errs.ErrNotFound
	}

	rr.logger.Infof("found %d expired reservations with readerID %s", len(reservations), readerID)

	return reservations, nil
}

func (rr *ReservationRepo) GetActiveByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	rr.logger.Infof("selecting active reservations with readerID: %s", readerID)

	query := fmt.Sprintf(`SELECT * FROM bs.reservation_view WHERE reader_id = $1 AND state != '%s' AND state != '%s'`, impl.ReservationClosed, impl.ReservationExpired)

	var reservations []*models.ReservationModel
	err := rr.db.SelectContext(ctx, &reservations, query, readerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting active reservations: %v", err)
		return nil, err
	}
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("active reservations with this readerID not found: %s", readerID)
		return nil, errs.ErrNotFound
	}

	rr.logger.Infof("found %d active reservations with readerID %s", len(reservations), readerID)

	return reservations, nil
}
