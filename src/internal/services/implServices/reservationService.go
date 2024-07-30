package implServices

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errsRepo"
	"BookSmart/internal/repositories/intfRepo"
	"BookSmart/internal/services/errsService"
	"BookSmart/internal/services/intfServices"
	"BookSmart/pkg/transact"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	ReservationIssued   = "Issued"
	ReservationExtended = "Extended"
	ReservationExpired  = "Expired"
	ReservationClosed   = "Closed"
)

const (
	ReservationIssuePeriodDays     = 14
	ReservationExtensionPeriodDays = 7
)

type ReservationService struct {
	reservationRepo    intfRepo.IReservationRepo
	bookRepo           intfRepo.IBookRepo
	readerRepo         intfRepo.IReaderRepo
	libCardRepo        intfRepo.ILibCardRepo
	transactionManager transact.ITransactionManager
	logger             *logrus.Entry
}

func NewReservationService(
	reservationRepo intfRepo.IReservationRepo,
	bookRepo intfRepo.IBookRepo,
	readerRepo intfRepo.IReaderRepo,
	libCardRepo intfRepo.ILibCardRepo,
	transactionManager transact.ITransactionManager,
	logger *logrus.Entry,
) intfServices.IReservationService {
	return &ReservationService{
		reservationRepo:    reservationRepo,
		bookRepo:           bookRepo,
		readerRepo:         readerRepo,
		libCardRepo:        libCardRepo,
		transactionManager: transactionManager,
		logger:             logger,
	}
}

func (rs *ReservationService) Create(ctx context.Context, readerID, bookID uuid.UUID) error {
	rs.logger.Info("starting reservation creation process")

	existingReader, err := rs.checkReader(ctx, readerID)
	if err != nil {
		return err
	}

	existingBook, err := rs.checkBook(ctx, bookID)
	if err != nil {
		return err
	}

	err = rs.checkAgeLimit(existingReader, existingBook)
	if err != nil {
		return err
	}

	err = rs.create(ctx, readerID, bookID)
	if err != nil {
		rs.logger.Errorf("error creating reservation: %v", err)
		return err
	}

	rs.logger.Info("reservation creation successful")

	return nil
}

func (rs *ReservationService) Update(ctx context.Context, reservation *models.ReservationModel) error {
	rs.logger.Info("attempting to update reservation")

	err := rs.checkValidLibCard(ctx, reservation.ReaderID)
	if err != nil {
		return err
	}
	err = rs.checkNoExpiredBooks(ctx, reservation.ReaderID)
	if err != nil {
		return err
	}

	err = rs.checkReservationState(reservation.State)
	if err != nil {
		return err
	}

	err = rs.checkBookRarityUpdate(ctx, reservation.BookID)
	if err != nil {
		return err
	}

	reservation.IssueDate = time.Now()
	reservation.ReturnDate = time.Now().AddDate(0, 0, ReservationExtensionPeriodDays)
	reservation.State = ReservationExtended

	rs.logger.Info("update reservation in repository")
	err = rs.reservationRepo.Update(ctx, reservation)
	if err != nil {
		rs.logger.Errorf("error updating reservation: %v", err)
		return err
	}

	rs.logger.Info("reservation update successful")

	return nil
}

func (rs *ReservationService) create(ctx context.Context, readerID, bookID uuid.UUID) error {
	return rs.transactionManager.Do(ctx, func(ctx context.Context) error {
		newReservation := &models.ReservationModel{
			ID:         uuid.New(),
			ReaderID:   readerID,
			BookID:     bookID,
			IssueDate:  time.Now(),
			ReturnDate: time.Now().AddDate(0, 0, ReservationIssuePeriodDays),
			State:      ReservationIssued,
		}

		rs.logger.Info("creating reservation in repository")

		err := rs.reservationRepo.Create(ctx, newReservation)
		if err != nil {
			rs.logger.Errorf("error creating reservation: %v", err)
			return err
		}

		existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
		if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
			rs.logger.Errorf("error checking book existence: %v", err)
			return err
		}

		existingBook.CopiesNumber -= 1

		rs.logger.Info("updating book copiesNumber in repository")

		err = rs.bookRepo.Update(ctx, existingBook)
		if err != nil {
			rs.logger.Errorf("error updating book: %v", err)
			return err
		}

		rs.logger.Info("successfully updated book copiesNumber")

		return nil
	})
}

func (rs *ReservationService) checkReader(ctx context.Context, readerID uuid.UUID) (*models.ReaderModel, error) {
	existingReader, err := rs.checkReaderExists(ctx, readerID)
	if err != nil {
		return nil, err
	}

	err = rs.checkNoExpiredBooks(ctx, readerID)
	if err != nil {
		return nil, err
	}

	err = rs.checkActiveReservationsLimit(ctx, readerID)
	if err != nil {
		return nil, err
	}

	err = rs.checkValidLibCard(ctx, readerID)
	if err != nil {
		return nil, err
	}

	rs.logger.Info("reader is valid")

	return existingReader, nil
}

func (rs *ReservationService) checkBook(ctx context.Context, bookID uuid.UUID) (*models.BookModel, error) {
	existingBook, err := rs.checkBookExists(ctx, bookID)
	if err != nil {
		return nil, err
	}

	err = rs.checkBookCopiesNumber(existingBook)
	if err != nil {
		return nil, err
	}

	err = rs.checkBookRarityCreate(existingBook)
	if err != nil {
		return nil, err
	}

	rs.logger.Info("book is valid")

	return existingBook, nil
}

func (rs *ReservationService) checkReaderExists(ctx context.Context, readerID uuid.UUID) (*models.ReaderModel, error) {
	existingReader, err := rs.readerRepo.GetByID(ctx, readerID)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		rs.logger.Errorf("error checking book existence: %v", err)
		return nil, err
	}
	if existingReader == nil {
		rs.logger.Warn("reader with this ID does not exist")
		return nil, errsService.ErrReaderDoesNotExists
	}

	rs.logger.Info("reader exists")

	return existingReader, nil
}

func (rs *ReservationService) checkNoExpiredBooks(ctx context.Context, readerID uuid.UUID) error {
	overdueBooks, err := rs.reservationRepo.GetExpiredByReaderID(ctx, readerID)
	if err != nil {
		rs.logger.Errorf("error checking expired book existence: %v", err)
		return err
	}

	if len(overdueBooks) > 0 {
		rs.logger.Warn("reader has expired books")
		return errsService.ErrReaderHasExpiredBooks
	}

	rs.logger.Info("reader has not expired books")

	return nil
}

func (rs *ReservationService) checkActiveReservationsLimit(ctx context.Context, readerID uuid.UUID) error {
	activeReservations, err := rs.reservationRepo.GetActiveByReaderID(ctx, readerID)
	if err != nil {
		rs.logger.Errorf("error checking active reservations: %v", err)
		return err
	}
	if len(activeReservations) >= MaxBooksPerReader {
		rs.logger.Warn("reader has reached the limit of active reservations")
		return errsService.ErrReservationsLimitExceeded
	}

	rs.logger.Info("reader has not reached the limit of active reservations")

	return nil
}

func (rs *ReservationService) checkValidLibCard(ctx context.Context, readerID uuid.UUID) error {
	libCard, err := rs.libCardRepo.GetByReaderID(ctx, readerID)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		rs.logger.Errorf("error checking libCard existence: %v", err)
		return err
	}
	if libCard == nil {
		rs.logger.Warn("reader does not have libCard")
		return errsService.ErrLibCardDoesNotExists
	}

	if !libCard.ActionStatus {
		rs.logger.Warn("reader has invalid libCard")
		return errsService.ErrLibCardIsInvalid
	}

	rs.logger.Info("reader has valid libCard")

	return nil
}

func (rs *ReservationService) checkBookExists(ctx context.Context, bookID uuid.UUID) (*models.BookModel, error) {
	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		rs.logger.Errorf("error checking book existence: %v", err)
		return nil, err
	}
	if existingBook == nil {
		rs.logger.Warn("book with this ID does not exist")
		return nil, errsService.ErrBookDoesNotExists
	}

	rs.logger.Info("book exists")

	return existingBook, nil
}

func (rs *ReservationService) checkBookCopiesNumber(book *models.BookModel) error {
	if book.CopiesNumber <= 0 {
		rs.logger.Warn("no copies of the book are available in the library")
		return errsService.ErrBookNoCopiesNum
	}

	rs.logger.Info("book has copies available")

	return nil
}

func (rs *ReservationService) checkBookRarityCreate(book *models.BookModel) error {
	if book.Rarity == BookRarityUnique {
		rs.logger.Warn("this book is unique and cannot be reserved")
		return errsService.ErrUniqueBookNotReserved
	}

	rs.logger.Info("book is not unique")

	return nil
}

func (rs *ReservationService) checkAgeLimit(reader *models.ReaderModel, book *models.BookModel) error {
	if reader.Age < book.AgeLimit {
		rs.logger.Warn("reader does not meet the age requirement for this book")
		return errsService.ErrReservationAgeLimit
	}

	rs.logger.Info("reader's age is appropriate")

	return nil
}

func (rs *ReservationService) checkBookRarityUpdate(ctx context.Context, bookID uuid.UUID) error {
	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil && !errors.Is(err, errsRepo.ErrNotFound) {
		rs.logger.Errorf("error checking book existence: %v", err)
		return err
	}

	if existingBook.Rarity == BookRarityRare || existingBook.Rarity == BookRarityUnique {
		rs.logger.Warn("rare and unique book cannot be renewed.")
		return errsService.ErrRareAndUniqueBookNotExtended
	}

	rs.logger.Info("book's rarity is common")

	return nil
}

func (rs *ReservationService) checkReservationState(reservationState string) error {
	if reservationState == ReservationClosed {
		rs.logger.Warn("this reservation is already closed")
		return errsService.ErrReservationIsAlreadyClosed
	}

	if reservationState == ReservationExpired {
		rs.logger.Warn("this reservation is already expired")
		return errsService.ErrReservationIsAlreadyExpired
	}

	if reservationState == ReservationExtended {
		rs.logger.Warn("this reservation is already extended")
		return errsService.ErrReservationIsAlreadyExtended
	}

	rs.logger.Info("reservation is only issued")

	return nil
}
