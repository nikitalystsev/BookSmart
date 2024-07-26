package postgres

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type ReaderRepo struct {
	db *sqlx.DB
}

func NewReaderRepo(db *sqlx.DB) intfRepo.IReaderRepo {
	return &ReaderRepo{db: db}
}

func (rr ReaderRepo) Create(ctx context.Context, reader *models.ReaderModel) error {
	//TODO implement me
	panic("implement me")
}

func (rr ReaderRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.ReaderModel, error) {
	//TODO implement me
	panic("implement me")
}

func (rr ReaderRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.ReaderModel, error) {
	//TODO implement me
	panic("implement me")
}

func (rr ReaderRepo) IsFavorite(ctx context.Context, readerID, bookID uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (rr ReaderRepo) AddToFavorites(ctx context.Context, readerID, bookID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (rr ReaderRepo) SaveRefreshToken(ctx context.Context, id uuid.UUID, token string, ttl time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (rr ReaderRepo) GetByRefreshToken(ctx context.Context, token string) (*models.ReaderModel, error) {
	//TODO implement me
	panic("implement me")
}
