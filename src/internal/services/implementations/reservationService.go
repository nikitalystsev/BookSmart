package implementations

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/repositories/interfaces"
	"BookSmart/internal/transactionManager/implementations"
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
	reservationRepo    interfaces.IReservationRepo
	bookRepo           interfaces.IBookRepo
	readerRepo         interfaces.IReaderRepo
	libCardRepo        interfaces.ILibCardRepo
	transactionManager implementations.TransactionManager
}

func NewReservationService(
	reservationRepo interfaces.IReservationRepo,
	bookRepo interfaces.IBookRepo,
	readerRepo interfaces.IReaderRepo,
	libCardRepo interfaces.ILibCardRepo,
	transactionManager implementations.TransactionManager,
) *ReservationService {
	return &ReservationService{
		reservationRepo:    reservationRepo,
		bookRepo:           bookRepo,
		readerRepo:         readerRepo,
		libCardRepo:        libCardRepo,
		transactionManager: transactionManager,
	}
}

// Create TODO делать в рамках одной транзакции (выполнено)
// TODO вынести проверки из транзакции
func (rs *ReservationService) Create(ctx context.Context, readerID, bookID uuid.UUID) error {
	return rs.transactionManager.WithTransaction(ctx, func(ctx context.Context) error {
		existingReader, err := rs.readerRepo.GetByID(ctx, readerID)
		if err != nil && !errors.Is(err, errs.ErrNotFound) {
			return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
		}
		if existingReader == nil {
			return fmt.Errorf("[!] ERROR! Reader with this ID does not exist")
		}

		existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
		if err != nil && !errors.Is(err, errs.ErrNotFound) {
			return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
		}
		if existingBook == nil {
			return fmt.Errorf("[!] ERROR! Book with this ID does not exist")
		}

		libCard, err := rs.libCardRepo.GetByReaderID(ctx, readerID)
		if err != nil {
			return fmt.Errorf("[!] ERROR! Error checking libCard existence: %v", err)
		}
		if libCard == nil || !libCard.ActionStatus {
			return fmt.Errorf("[!] ERROR! Reader does not have a valid library card")
		}

		overdueBooks, err := rs.reservationRepo.GetOverdueByReaderID(ctx, readerID)
		if err != nil {
			return fmt.Errorf("[!] ERROR! Error checking overdue books: %v", err)
		}
		if len(overdueBooks) > 0 {
			return fmt.Errorf("[!] ERROR! Reader has overdue books")
		}

		activeReservations, err := rs.reservationRepo.GetActiveByReaderID(ctx, readerID)
		if err != nil {
			return fmt.Errorf("[!] ERROR! Error checking active reservations: %v", err)
		}
		if len(activeReservations) >= MaxBooksPerReader {
			return fmt.Errorf("[!] ERROR! Reader has reached the limit of active reservations")
		}

		if existingBook.CopiesNumber <= 0 {
			return fmt.Errorf("[!] ERROR! No copies of the book are available in the library")
		}

		if existingReader.Age < existingBook.AgeLimit {
			return fmt.Errorf("[!] ERROR! Reader does not meet the age requirement for this book")
		}

		if existingBook.Rarity == BookRarityUnique {
			return fmt.Errorf("[!] ERROR! This book is unique and cannot be reserved")
		}

		newReservation := &models.ReservationModel{
			ID:         uuid.New(),
			ReaderID:   readerID,
			BookID:     bookID,
			IssueDate:  time.Now(),
			ReturnDate: time.Now().AddDate(0, 0, ReservationIssuePeriodDays),
			State:      ReservationIssued,
		}

		err = rs.reservationRepo.Create(ctx, newReservation)
		if err != nil {
			return fmt.Errorf("[!] ERROR! Error creating reservation: %v", err)
		}

		existingBook.CopiesNumber -= 1

		err = rs.bookRepo.Update(ctx, existingBook)
		if err != nil {
			return fmt.Errorf("[!] ERROR! Error updating book availability: %v", err)
		}

		return nil
	})
}

func (rs *ReservationService) Update(ctx context.Context, reservation *models.ReservationModel) error {
	libCard, err := rs.libCardRepo.GetByReaderID(ctx, reservation.ReaderID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking libCard existence: %v", err)
	}
	if libCard == nil || !libCard.ActionStatus {
		return fmt.Errorf("[!] ERROR! Reader does not have a valid library card")
	}

	overdueBooks, err := rs.reservationRepo.GetOverdueByReaderID(ctx, reservation.ReaderID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking overdue books: %v", err)
	}
	if len(overdueBooks) > 0 {
		return fmt.Errorf("[!] ERROR! Reader has overdue books")
	}

	if reservation.State == ReservationClosed {
		return fmt.Errorf("[!] ERROR! This reservation is already closed")
	}

	if reservation.State == ReservationOverdue {
		return fmt.Errorf("[!] ERROR! This reservation is past its return date")
	}

	if reservation.State == ReservationExtended {
		return fmt.Errorf("[!] ERROR! This reservation has already been extended")
	}

	existingBook, err := rs.bookRepo.GetByID(ctx, reservation.BookID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}

	if existingBook.Rarity != BookRarityCommon {
		return fmt.Errorf("[!] ERROR! This book is not renewed")
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
