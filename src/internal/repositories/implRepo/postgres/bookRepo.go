package postgres

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/intfRepo"
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BookRepo struct {
	db *sqlx.DB
}

func NewBookRepo(db *sqlx.DB) intfRepo.IBookRepo {
	return &BookRepo{db: db}
}

func (br BookRepo) Create(ctx context.Context, book *models.BookModel) error {
	//TODO implement me
	panic("implement me")
}

func (br BookRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.BookModel, error) {
	//TODO implement me
	panic("implement me")
}

func (br BookRepo) GetByTitle(ctx context.Context, title string) (*models.BookModel, error) {
	//TODO implement me
	panic("implement me")
}

func (br BookRepo) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (br BookRepo) Update(ctx context.Context, book *models.BookModel) error {
	//TODO implement me
	panic("implement me")
}

func (br BookRepo) GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error) {
	//TODO implement me
	panic("implement me")
}
