package postgres

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReservationRepo struct {
	db *sqlx.DB
}

func (rr ReservationRepo) Create(ctx context.Context, reservation *models.ReservationModel) error {
	//TODO implement me
	panic("implement me")
}

func (rr ReservationRepo) GetByReaderAndBook(ctx context.Context, readerID, bookID uuid.UUID) (*models.ReservationModel, error) {
	//TODO implement me
	panic("implement me")
}

func (rr ReservationRepo) GetByID(ctx context.Context, reservationID uuid.UUID) (*models.ReservationModel, error) {
	//TODO implement me
	panic("implement me")
}

func (rr ReservationRepo) Update(ctx context.Context, reservation *models.ReservationModel) error {
	//TODO implement me
	panic("implement me")
}

func (rr ReservationRepo) GetOverdueByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	//TODO implement me
	panic("implement me")
}

func (rr ReservationRepo) GetActiveByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	//TODO implement me
	panic("implement me")
}

func NewReservationRepo(db *sqlx.DB) intfRepo.IReservationRepo {
	return &ReservationRepo{db: db}
}
