package impl

import (
	"BookSmart-services/core/models"
	"BookSmart-services/intfRepo"
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type LibCardRepo struct {
	db     *mongo.Collection
	logger *logrus.Entry
}

func NewLibCardRepo(db *mongo.Database, logger *logrus.Entry) intfRepo.ILibCardRepo {
	return &LibCardRepo{db: db.Collection("lib_card"), logger: logger}
}

func (l *LibCardRepo) Create(ctx context.Context, libCard *models.LibCardModel) error {
	//TODO implement me
	panic("implement me")
}

func (l *LibCardRepo) GetByReaderID(ctx context.Context, readerID uuid.UUID) (*models.LibCardModel, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LibCardRepo) GetByNum(ctx context.Context, libCardNum string) (*models.LibCardModel, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LibCardRepo) Update(ctx context.Context, libCard *models.LibCardModel) error {
	//TODO implement me
	panic("implement me")
}
