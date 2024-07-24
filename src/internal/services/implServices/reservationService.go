package implServices

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/repositories/intfRepo"
	"BookSmart/pkg/transact"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	ReservationIssued   = "Выдана"
	ReservationExtended = "Продлена"
	ReservationOverdue  = "Просрочена"
	ReservationClosed   = "Закрыта"
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
}

func NewReservationService(
	reservationRepo intfRepo.IReservationRepo,
	bookRepo intfRepo.IBookRepo,
	readerRepo intfRepo.IReaderRepo,
	libCardRepo intfRepo.ILibCardRepo,
	transactionManager transact.ITransactionManager,
) *ReservationService {
	return &ReservationService{
		reservationRepo:    reservationRepo,
		bookRepo:           bookRepo,
		readerRepo:         readerRepo,
		libCardRepo:        libCardRepo,
		transactionManager: transactionManager,
	}
}

func (rs *ReservationService) Create(ctx context.Context, readerID, bookID uuid.UUID) error {
	err := rs.checkReaderExists(ctx, readerID)
	if err != nil {
		return err
	}

	err = rs.checkBookExists(ctx, bookID)
	if err != nil {
		return err
	}

	err = rs.checkValidLibCard(ctx, readerID)
	if err != nil {
		return err
	}

	err = rs.checkNoOverdueBooks(ctx, readerID)
	if err != nil {
		return err
	}

	err = rs.checkActiveReservationsLimit(ctx, readerID)
	if err != nil {
		return err
	}

	err = rs.checkBookCopiesNumber(ctx, bookID)
	if err != nil {
		return err
	}

	err = rs.checkAgeLimit(ctx, readerID, bookID)
	if err != nil {
		return err
	}

	err = rs.checkBookRarityCreate(ctx, bookID)
	if err != nil {
		return err
	}
	err = rs.create(ctx, readerID, bookID)
	if err != nil {
		return err
	}
	return nil
}

func (rs *ReservationService) Update(ctx context.Context, reservation *models.ReservationModel) error {
	err := rs.checkValidLibCard(ctx, reservation.ReaderID)
	if err != nil {
		return err
	}
	err = rs.checkNoOverdueBooks(ctx, reservation.BookID)
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

	err = rs.reservationRepo.Update(ctx, reservation)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error creating reservation: %v", err)
	}

	return nil
}

func (rs *ReservationService) create(ctx context.Context, readerID, bookID uuid.UUID) error {
	return rs.transactionManager.WithTransaction(ctx, func(ctx context.Context) error {
		newReservation := &models.ReservationModel{
			ID:         uuid.New(),
			ReaderID:   readerID,
			BookID:     bookID,
			IssueDate:  time.Now(),
			ReturnDate: time.Now().AddDate(0, 0, ReservationIssuePeriodDays),
			State:      ReservationIssued,
		}

		err := rs.reservationRepo.Create(ctx, newReservation)
		if err != nil {
			return fmt.Errorf("[!] ERROR! Error creating reservation: %v", err)
		}

		existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
		if err != nil {
			return fmt.Errorf("[!] ERROR! Error retrieving book: %v", err)
		}

		existingBook.CopiesNumber -= 1

		err = rs.bookRepo.Update(ctx, existingBook)
		if err != nil {
			return fmt.Errorf("[!] ERROR! Error updating book availability: %v", err)
		}

		return nil
	})
}

func (rs *ReservationService) checkReaderExistence(ctx context.Context, readerID uuid.UUID) error {
	existingReader, err := rs.readerRepo.GetByID(ctx, readerID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}
	if existingReader == nil {
		return fmt.Errorf("[!] ERROR! Reader with this ID does not exist")
	}

	return nil
}

func (rs *ReservationService) checkReaderExists(ctx context.Context, readerID uuid.UUID) error {
	existingReader, err := rs.readerRepo.GetByID(ctx, readerID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}
	if existingReader == nil {
		return fmt.Errorf("[!] ERROR! Reader with this ID does not exist")
	}
	return nil
}

func (rs *ReservationService) checkBookExists(ctx context.Context, bookID uuid.UUID) error {
	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}
	if existingBook == nil {
		return fmt.Errorf("[!] ERROR! Book with this ID does not exist")
	}
	return nil
}

func (rs *ReservationService) checkValidLibCard(ctx context.Context, readerID uuid.UUID) error {
	libCard, err := rs.libCardRepo.GetByReaderID(ctx, readerID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking libCard existence: %v", err)
	}
	if libCard == nil || !libCard.ActionStatus {
		return fmt.Errorf("[!] ERROR! Reader does not have a valid library card")
	}
	return nil
}

func (rs *ReservationService) checkNoOverdueBooks(ctx context.Context, readerID uuid.UUID) error {
	overdueBooks, err := rs.reservationRepo.GetOverdueByReaderID(ctx, readerID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking overdue books: %v", err)
	}
	if len(overdueBooks) > 0 {
		return fmt.Errorf("[!] ERROR! Reader has overdue books")
	}
	return nil
}

func (rs *ReservationService) checkActiveReservationsLimit(ctx context.Context, readerID uuid.UUID) error {
	activeReservations, err := rs.reservationRepo.GetActiveByReaderID(ctx, readerID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking active reservations: %v", err)
	}
	if len(activeReservations) >= MaxBooksPerReader {
		return fmt.Errorf("[!] ERROR! Reader has reached the limit of active reservations")
	}
	return nil
}

func (rs *ReservationService) checkBookCopiesNumber(ctx context.Context, bookID uuid.UUID) error {
	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}
	if existingBook.CopiesNumber <= 0 {
		return fmt.Errorf("[!] ERROR! No copies of the book are available in the library")
	}
	return nil
}

func (rs *ReservationService) checkAgeLimit(ctx context.Context, readerID, bookID uuid.UUID) error {
	existingReader, err := rs.readerRepo.GetByID(ctx, readerID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error retrieving reader: %v", err)
	}

	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error retrieving book: %v", err)
	}

	if existingReader.Age < existingBook.AgeLimit {
		return fmt.Errorf("[!] ERROR! Reader does not meet the age requirement for this book")
	}
	return nil
}

func (rs *ReservationService) checkBookRarityCreate(ctx context.Context, bookID uuid.UUID) error {
	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}
	if existingBook.Rarity == BookRarityUnique {
		return fmt.Errorf("[!] ERROR! This book is unique and cannot be reserved")
	}
	return nil
}

func (rs *ReservationService) checkBookRarityUpdate(ctx context.Context, bookID uuid.UUID) error {
	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}

	if existingBook.Rarity == BookRarityRare || existingBook.Rarity == BookRarityUnique {
		return fmt.Errorf("[!] ERROR! This book is not renewed")
	}

	return nil
}

func (rs *ReservationService) checkReservationState(reservationState string) error {
	if reservationState == ReservationClosed {
		return fmt.Errorf("[!] ERROR! This reservation is already closed")
	}

	if reservationState == ReservationOverdue {
		return fmt.Errorf("[!] ERROR! This reservation is past its return date")
	}

	if reservationState == ReservationExtended {
		return fmt.Errorf("[!] ERROR! This reservation has already been extended")
	}

	return nil
}
