package intfServices

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"context"
)

type IReaderService interface {
	SignUp(ctx context.Context, reader *models.ReaderModel) error
	SignIn(ctx context.Context, reader *dto.ReaderLoginDTO) error
}
