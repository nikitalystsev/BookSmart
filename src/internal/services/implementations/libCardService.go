package implementations

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/repositories/interfaces"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/big"
	"time"
)

const (
	libCardNumLength = 13
	charset          = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type LibCardService struct {
	libCardRepo interfaces.ILibCardRepo
}

func NewLibCardService(libCardRepo interfaces.ILibCardRepo) *LibCardService {
	return &LibCardService{libCardRepo: libCardRepo}
}

func (lcs *LibCardService) Create(ctx context.Context, readerID uuid.UUID) error {
	existingLibCard, err := lcs.libCardRepo.GetByReaderID(ctx, readerID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking libCard existence: %v", err)
	}

	if existingLibCard != nil {
		return fmt.Errorf("[!] ERROR! User with ID %v already has a library card", readerID)
	}

	libCardNum, _ := lcs.generateLibCardNum()

	newLibCard := &models.LibCardModel{
		ID:           uuid.New(),
		ReaderID:     readerID,
		LibCardNum:   libCardNum,
		Validity:     365, // Срок действия 1 год (365 дней)
		IssueDate:    time.Now(),
		ActionStatus: true,
	}

	err = lcs.libCardRepo.Create(ctx, newLibCard)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error creating libCard: %v", err)
	}

	return nil
}

func (lcs *LibCardService) Update(libCard *models.LibCardModel) error {
	// логика обновления
	return nil
}

func (lcs *LibCardService) IsValidLibCard(libCard *models.LibCardModel) bool {
	if !libCard.ActionStatus {
		return false
	}

	expiryDate := libCard.IssueDate.AddDate(0, 0, libCard.Validity)

	return time.Now().Before(expiryDate)
}

func (lcs *LibCardService) generateLibCardNum() (string, error) {
	result := make([]byte, libCardNumLength)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("[!] ERROR! Error generating library card number: %v", err)
		}

		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
