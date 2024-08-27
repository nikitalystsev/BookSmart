package impl

import (
	"BookSmart-services/core/models"
	"BookSmart-services/errs"
	"BookSmart-services/impl"
	"BookSmart-services/intfRepo"
	"context"
	"database/sql"
	"errors"
	"fmt"
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

type ReservationRepo struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
	logger *logrus.Entry
}

func NewReservationRepo(db *sqlx.DB, logger *logrus.Entry) intfRepo.IReservationRepo {
	return &ReservationRepo{db: db, getter: trmsqlx.DefaultCtxGetter, logger: logger}
}

func (rr *ReservationRepo) Create(ctx context.Context, reservation *models.ReservationModel) error {
	rr.logger.Infof("inserting reservation with ID: %s", reservation.ID)

	query := `insert into bs.reservation values ($1, $2, $3, $4, $5, $6)`

	result, err := rr.getter.DefaultTrOrDB(ctx, rr.db).ExecContext(ctx, query, reservation.ID, reservation.ReaderID, reservation.BookID,
		reservation.IssueDate, reservation.ReturnDate, reservation.State)
	if err != nil {
		rr.logger.Errorf("error inserting reservation: %v", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		rr.logger.Errorf("error inserting reservation: %v", err)
		return err
	}
	if rows != 1 {
		rr.logger.Errorf("error inserting reservation: %d rows affected", rows)
		return errors.New("reservationRepo.Create: expected 1 row affected")
	}

	rr.logger.Infof("inserted reservation with ID: %s", reservation.ID)

	return nil
}

func (rr *ReservationRepo) GetByReaderAndBook(ctx context.Context, readerID, bookID uuid.UUID) (*models.ReservationModel, error) {
	rr.logger.Infof("selecting reservation with readerID и bookID: %s и %s", readerID, bookID)

	query := `select * from bs.reservation_view where reader_id = $1 and book_id = $2`

	var reservation models.ReservationModel
	err := rr.getter.DefaultTrOrDB(ctx, rr.db).GetContext(ctx, &reservation, query, readerID, bookID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting reservation: %v", err)
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("reservation with this readerID и bookID not found: %s и %s", readerID, bookID)
		return nil, errs.ErrReservationDoesNotExists
	}

	rr.logger.Infof("selected reservation with readerID и bookID: %s и %s", readerID, bookID)

	return &reservation, nil
}

func (rr *ReservationRepo) GetByID(ctx context.Context, ID uuid.UUID) (*models.ReservationModel, error) {
	rr.logger.Infof("selecting reservation with ID: %s", ID)

	query := `select * from bs.reservation_view where id = $1`

	var reservation models.ReservationModel
	err := rr.getter.DefaultTrOrDB(ctx, rr.db).GetContext(ctx, &reservation, query, ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting reservation with ID: %v", err)
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("reservation with this ID not found: %s", ID)
		return nil, errs.ErrReservationDoesNotExists
	}

	rr.logger.Infof("selected reservation with ID: %s", ID)

	return &reservation, nil
}

// GetByBookID TODO добавить в схемы (протестировано)
func (rr *ReservationRepo) GetByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.ReservationModel, error) {
	rr.logger.Infof("selecting reservation with bookID: %s", bookID)

	query := fmt.Sprintf(`select * from bs.reservation_view where book_id = $1 and state != '%s'`, impl.ReservationClosed)

	var reservations []*models.ReservationModel
	err := rr.getter.DefaultTrOrDB(ctx, rr.db).SelectContext(ctx, &reservations, query, bookID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting reservation with ID: %v", err)
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		rr.logger.Warnf("reservation with this bookID not found: %s", bookID)
		return nil, errs.ErrReservationDoesNotExists
	}

	rr.logger.Infof("selected reservation with bookID: %s", bookID)

	return reservations, nil
}

func (rr *ReservationRepo) Update(ctx context.Context, reservation *models.ReservationModel) error {
	rr.logger.Infof("updating reservation with ID: %s", reservation.ID)

	query := `update bs.reservation set issue_date = $1, return_date = $2, state = $3 where id = $4`

	result, err := rr.getter.DefaultTrOrDB(ctx, rr.db).ExecContext(ctx, query, reservation.IssueDate, reservation.ReturnDate, reservation.State, reservation.ID)
	if err != nil {
		rr.logger.Errorf("error updating reservation with ID: %v", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		rr.logger.Errorf("error updating reservation with ID: %v", err)
		return err
	}
	if rows != 1 {
		rr.logger.Errorf("error updating reservation: %d rows affected", rows)
		return errors.New("reservationRepo.Update: expected 1 row affected")
	}

	rr.logger.Infof("updated reservation with ID: %s", reservation.ID)

	return nil
}

func (rr *ReservationRepo) GetExpiredByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	rr.logger.Infof("selecting expired reservations with readerID: %s", readerID)

	query := `select * from bs.reservation_view where reader_id = $1 and return_date < $2`

	var reservations []*models.ReservationModel
	err := rr.getter.DefaultTrOrDB(ctx, rr.db).SelectContext(ctx, &reservations, query, readerID, time.Now())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting expired reservations: %v", err)
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) || len(reservations) == 0 {
		rr.logger.Warnf("expired reservations with this readerID not found: %s", readerID)
		return nil, errs.ErrReservationDoesNotExists
	}

	rr.logger.Infof("found %d expired reservations with readerID %s", len(reservations), readerID)

	return reservations, nil
}

func (rr *ReservationRepo) GetActiveByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	rr.logger.Infof("selecting active reservations with readerID: %s", readerID)

	query := fmt.Sprintf(`select * from bs.reservation_view where reader_id = $1 and state != '%s' and state != '%s'`, impl.ReservationClosed, impl.ReservationExpired)

	var reservations []*models.ReservationModel
	err := rr.getter.DefaultTrOrDB(ctx, rr.db).SelectContext(ctx, &reservations, query, readerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rr.logger.Errorf("error selecting active reservations: %v", err)
		return nil, err
	}
	if errors.Is(err, sql.ErrNoRows) || len(reservations) == 0 {
		rr.logger.Warnf("active reservations with this readerID not found: %s", readerID)
		return nil, errs.ErrReservationDoesNotExists
	}

	rr.logger.Infof("found %d active reservations with readerID %s", len(reservations), readerID)

	return reservations, nil
}
