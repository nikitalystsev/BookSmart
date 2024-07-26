package postgres

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LibCardRepo struct {
	db *sqlx.DB
}

func NewLibCardRepo(db *sqlx.DB) intfRepo.ILibCardRepo {
	return &LibCardRepo{db: db}
}

func (lcr LibCardRepo) Create(ctx context.Context, libCard *models.LibCardModel) error {
	//TODO implement me
	panic("implement me")
}

func (lcr LibCardRepo) GetByReaderID(ctx context.Context, readerID uuid.UUID) (*models.LibCardModel, error) {
	//TODO implement me
	panic("implement me")
}

func (lcr LibCardRepo) GetByNum(ctx context.Context, libCardNum string) (*models.LibCardModel, error) {
	//TODO implement me
	panic("implement me")
}

func (lcr LibCardRepo) Update(ctx context.Context, libCard *models.LibCardModel) error {
	//TODO implement me
	panic("implement me")
}
