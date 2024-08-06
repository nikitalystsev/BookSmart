package impl

import (
	"BookSmart-services/core/dto"
	"BookSmart-services/core/models"
	"BookSmart-services/intfRepo"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepo struct {
	db     *mongo.Collection
	logger *logrus.Entry
}

func NewBookRepo(db *mongo.Database, logger *logrus.Entry) intfRepo.IBookRepo {
	return &BookRepo{db: db.Collection("book"), logger: logger}
}

func (b *BookRepo) Create(ctx context.Context, book *models.BookModel) error {
	//TODO implement me
	panic("implement me")
}

func (b *BookRepo) GetByID(ctx context.Context, ID uuid.UUID) (*models.BookModel, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookRepo) GetByTitle(ctx context.Context, title string) (*models.BookModel, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookRepo) Delete(ctx context.Context, ID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (b *BookRepo) Update(ctx context.Context, book *models.BookModel) error {
	//TODO implement me
	panic("implement me")
}

func (b *BookRepo) GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error) {
	//TODO implement me
	panic("implement me")
}
