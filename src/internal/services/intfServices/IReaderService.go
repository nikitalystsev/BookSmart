package intfServices

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"context"
)

type IReaderService interface {
	Register(ctx context.Context, reader *models.ReaderModel) error
	Login(ctx context.Context, reader *dto.ReaderLoginDTO) error
}
