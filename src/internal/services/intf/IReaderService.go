package intf

import (
	"BookSmart-services/core/dto"
	"BookSmart-services/core/models"
	"context"
	"github.com/google/uuid"
)

type IReaderService interface {
	SignUp(ctx context.Context, reader *models.ReaderModel) error
	SignIn(ctx context.Context, reader *dto.ReaderSignInDTO) (Tokens, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.ReaderModel, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	AddToFavorites(ctx context.Context, readerID, bookID uuid.UUID) error
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
