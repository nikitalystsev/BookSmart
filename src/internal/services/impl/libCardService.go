package impl

import (
	errsRepo "BookSmart-repositories/errs"
	intfRepo "BookSmart-repositories/intf"
	"BookSmart-services/errs"
	"BookSmart-services/intf"
	"BookSmart-services/models"
	"context"
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	logger      *logrus.Entry
}

func NewLibCardService(libCardRepo intfRepo.ILibCardRepo, logger *logrus.Entry) intf.ILibCardService {
	return &LibCardService{libCardRepo: libCardRepo, logger: logger}
}

func (lcs *LibCardService) Create(ctx context.Context, readerID uuid.UUID) error {
	lcs.logger.Info("starting libCard creation process")

	existingLibCard, err := lcs.libCardRepo.GetByReaderID(ctx, readerID)

	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		lcs.logger.Errorf("error checking libCard existence: %v", err)
		return err
	}

	if existingLibCard != nil {
		lcs.logger.Warnf("User with ID %v already has a library card", readerID)
		return errs.ErrLibCardAlreadyExist
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

	lcs.logger.Infof("creating libCard in repository: %+v", newLibCard)

	err = lcs.libCardRepo.Create(ctx, newLibCard)
	if err != nil {
		lcs.logger.Errorf("error creating libCard: %v", err)
		return err
	}

	lcs.logger.Info("libCard creation successful")

	return nil
}

func (lcs *LibCardService) Update(ctx context.Context, libCard *models.LibCardModel) error {
	if libCard == nil {
		lcs.logger.Warn("libCard object is nil")
		return errs.ErrLibCardObjectIsNil
	}

	lcs.logger.Infof("attempting to update libCard with ID: %s", libCard.ID)

	existingLibCard, err := lcs.libCardRepo.GetByNum(ctx, libCard.LibCardNum)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		lcs.logger.Errorf("error checking libCard existence: %v", err)
		return err
	}

	if existingLibCard == nil {
		lcs.logger.Warn("libCard with this Nun does not exist")
		return errs.ErrLibCardDoesNotExists
	}

	if lcs.isValidLibCard(existingLibCard) {
		lcs.logger.Warn("libCard with this Nun is already valid")
		return errs.ErrLibCardIsValid
	}

	libCard.IssueDate = time.Now()
	libCard.ActionStatus = true

	err = lcs.libCardRepo.Update(ctx, libCard)
	if err != nil {
		lcs.logger.Errorf("error updating libCard: %v", err)
		return err
	}

	lcs.logger.Infof("successfully updated book with ID: %s", libCard.ID)

	return nil
}

// GetByReaderID TODO добавить метод на схему
func (lcs *LibCardService) GetByReaderID(ctx context.Context, readerID uuid.UUID) (*models.LibCardModel, error) {
	lcs.logger.Infof("attempting to get libCard by readerID: %s", readerID)

	libCard, err := lcs.libCardRepo.GetByReaderID(ctx, readerID)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		lcs.logger.Errorf("error checking libCard existence: %v", err)
		return nil, err
	}

	if libCard == nil {
		lcs.logger.Warn("reader has no library card")
		return nil, errs.ErrLibCardDoesNotExists
	}

	lcs.logger.Infof("successfully getting libCard by readerID: %s", readerID)

	return libCard, nil
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