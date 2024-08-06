package impl

import (
	"BookSmart-services/core/models"
	"BookSmart-services/intfRepo"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ReaderRepo struct {
	db     *mongo.Collection
	logger *logrus.Entry
}

func NewReaderRepo(db *mongo.Collection, logger *logrus.Entry) intfRepo.IReaderRepo {
	return &ReaderRepo{db: db, logger: logger}
}

func (r *ReaderRepo) Create(ctx context.Context, reader *models.ReaderModel) error {
	//TODO implement me
	panic("implement me")
}

func (r *ReaderRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.ReaderModel, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReaderRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ReaderModel, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReaderRepo) IsFavorite(ctx context.Context, readerID, bookID uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *ReaderRepo) AddToFavorites(ctx context.Context, readerID, bookID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *ReaderRepo) SaveRefreshToken(ctx context.Context, id uuid.UUID, token string, ttl time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (r *ReaderRepo) GetByRefreshToken(ctx context.Context, token string) (*models.ReaderModel, error) {
	//TODO implement me
	panic("implement me")
}
