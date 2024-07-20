package impl

import (
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/repositories/interfaces"
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
	reservationRepo interfaces.IReservationRepo
	bookRepo        interfaces.IBookRepo
	readerRepo      interfaces.IReaderRepo
}

func NewReservationService(
	reservationRepo interfaces.IReservationRepo,
	bookRepo interfaces.IBookRepo,
	readerRepo interfaces.IReaderRepo,
) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		bookRepo:        bookRepo,
		readerRepo:      readerRepo,
	}
}

func (rs *ReservationService) Create(ctx context.Context, readerID, bookID uuid.UUID) error {

	exitingReader, err := rs.readerRepo.GetByID(ctx, readerID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}
	if exitingReader == nil {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: reader does not exist")
	}

	exitingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}
	if exitingBook == nil {
		return fmt.Errorf("[!] ERROR! Error checking book existence: book does not exist")
	}

	existingReservation, err := rs.reservationRepo.GetByReaderAndBook(ctx, readerID, bookID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking existing reservation: %v", err)
	}
	if existingReservation != nil {
		return fmt.Errorf("[!] ERROR! Book with ID %v is already reserved by reader with ID %v", bookID, readerID)
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

	return nil
}

func (rs *ReservationService) Update(ctx context.Context, reservation *models.ReservationModel) error {
	existingReservation, err := rs.reservationRepo.GetByID(ctx, reservation.ID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking reservation existence: %v", err)
	}

	if existingReservation == nil {
		return fmt.Errorf("[!] ERROR! Reservation with ID %v not found", reservation.ID)
	}

	err = rs.reservationRepo.Update(ctx, existingReservation)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error updating reservation: %v", err)
	}

	return nil
}
