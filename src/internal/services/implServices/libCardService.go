package implServices

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/repositories/intfRepo"
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

	LibCardValidityPeriod = 365
)

type LibCardService struct {
	libCardRepo intfRepo.ILibCardRepo
}

func NewLibCardService(libCardRepo intfRepo.ILibCardRepo) *LibCardService {
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

	libCardNum := lcs.generateLibCardNum()

	newLibCard := &models.LibCardModel{
		ID:           uuid.New(),
		ReaderID:     readerID,
		LibCardNum:   libCardNum,
		Validity:     LibCardValidityPeriod, // Срок действия 1 год (365 дней)
		IssueDate:    time.Now(),
		ActionStatus: true,
	}

	err = lcs.libCardRepo.Create(ctx, newLibCard)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error creating libCard: %v", err)
	}

	return nil
}

func (lcs *LibCardService) Update(ctx context.Context, libCard *models.LibCardModel) error {
	existingLibCard, err := lcs.libCardRepo.GetByNum(ctx, libCard.LibCardNum)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking libCard existence: %v", err)
	}

	if existingLibCard == nil {
		return fmt.Errorf("[!] ERROR! libCard with ID %v does not exist", libCard.LibCardNum)
	}

	if lcs.isValidLibCard(existingLibCard) {
		return fmt.Errorf("[!] ERROR! libCard with ID %v is valid", libCard.LibCardNum)
	}

	libCard.IssueDate = time.Now()

	err = lcs.libCardRepo.Update(ctx, libCard)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error updating libCard: %v", err)
	}

	return nil
}

func (lcs *LibCardService) isValidLibCard(libCard *models.LibCardModel) bool {
	if !libCard.ActionStatus {
		return false
	}

	expiryDate := libCard.IssueDate.AddDate(0, 0, libCard.Validity)

	return time.Now().Before(expiryDate)
}

func (lcs *LibCardService) generateLibCardNum() string {
	result := make([]byte, libCardNumLength)

	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}
