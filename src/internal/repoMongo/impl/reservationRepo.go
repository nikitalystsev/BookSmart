package impl

import (
	"BookSmart-services/core/models"
	"BookSmart-services/intfRepo"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReservationRepo struct {
	db     *mongo.Collection
	logger *logrus.Entry
}

func NewReservationRepo(db *mongo.Collection, logger *logrus.Entry) intfRepo.IReservationRepo {
	return &ReservationRepo{db: db, logger: logger}
}

func (r ReservationRepo) Create(ctx context.Context, reservation *models.ReservationModel) error {
	//TODO implement me
	panic("implement me")
}

func (r ReservationRepo) GetByReaderAndBook(ctx context.Context, readerID, bookID uuid.UUID) (*models.ReservationModel, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReservationRepo) GetByID(ctx context.Context, ID uuid.UUID) (*models.ReservationModel, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReservationRepo) Update(ctx context.Context, reservation *models.ReservationModel) error {
	//TODO implement me
	panic("implement me")
}

func (r ReservationRepo) GetExpiredByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReservationRepo) GetActiveByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error) {
	//TODO implement me
	panic("implement me")
}
