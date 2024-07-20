package impl

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/repositories/errs"
	"BookSmart/internal/repositories/interfaces"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const MaxBooksPerReader = 5

type ReaderService struct {
	ReaderRepo      interfaces.IReaderRepo
	ReservationRepo interfaces.IReservationRepo
	BookRepo        interfaces.IBookRepo
	LibCardRepo     interfaces.ILibCardRepo
}

func CreateNewReaderService(
	readerRepo interfaces.IReaderRepo,
	reservationRepo interfaces.IReservationRepo,
	bookRepo interfaces.IBookRepo,
	libCardRepo interfaces.ILibCardRepo,
) *ReaderService {
	return &ReaderService{
		ReaderRepo:      readerRepo,
		ReservationRepo: reservationRepo,
		BookRepo:        bookRepo,
		LibCardRepo:     libCardRepo,
	}
}

func (rs *ReaderService) Register(ctx context.Context, reader *models.ReaderModel) error {
	existingReader, err := rs.ReaderRepo.GetByPhoneNumber(ctx, reader.PhoneNumber)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}

	if existingReader != nil {
		return errors.New("[!] ERROR! Reader with this phoneNumbers already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reader.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error hashing password: %v", err)
	}

	reader.Password = string(hashedPassword)

	err = rs.ReaderRepo.Create(ctx, reader)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error creating reader: %v", err)
	}

	return nil
}

func (rs *ReaderService) Login(ctx context.Context, reader *dto.ReaderLoginDTO) error {
	exitingReader, err := rs.ReaderRepo.GetByPhoneNumber(ctx, reader.PhoneNumber)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}

	if exitingReader == nil {
		return fmt.Errorf("[!] ERROR! Reader with this phoneNumbers does not exist")
	}

	err = bcrypt.CompareHashAndPassword([]byte(exitingReader.Password), []byte(reader.Password))
	if err != nil {
		return fmt.Errorf("[!] ERROR! Wrong password")
	}

	return nil
}

func (rs *ReaderService) ReserveBook(ctx context.Context, readerID, bookID uuid.UUID) error {
	existingReader, err := rs.ReaderRepo.GetByID(ctx, readerID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking reader existence: %v", err)
	}
	if existingReader == nil {
		return fmt.Errorf("[!] ERROR! Reader with this ID does not exist")
	}

	existingBook, err := rs.BookRepo.GetByID(ctx, bookID)
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}
	if existingBook == nil {
		return fmt.Errorf("[!] ERROR! Book with this ID does not exist")
	}

	libCard, err := rs.LibCardRepo.GetByReaderID(ctx, readerID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking libCard existence: %v", err)
	}
	if libCard == nil || !rs.isValidLibCard(libCard) {
		return fmt.Errorf("[!] ERROR! Reader does not have a valid library card")
	}

	overdueBooks, err := rs.ReservationRepo.GetOverdueByReaderID(ctx, readerID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking overdue books: %v", err)
	}
	if len(overdueBooks) > 0 {
		return fmt.Errorf("[!] ERROR! Reader has overdue books")
	}

	activeReservations, err := rs.ReservationRepo.GetActiveByReaderID(ctx, readerID)
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

	err = rs.ReservationRepo.Create(ctx, newReservation)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error creating reservation: %v", err)
	}

	// TODO перенести в репозиторий, наверное
	existingBook.CopiesNumber -= 1

	err = rs.BookRepo.Update(ctx, existingBook)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error updating book availability: %v", err)
	}

	return nil
}

func (rs *ReaderService) ExtendBook(ctx context.Context, reservation *models.ReservationModel) error {
	libCard, err := rs.LibCardRepo.GetByReaderID(ctx, reservation.ReaderID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking libCard existence: %v", err)
	}
	if libCard == nil || !rs.isValidLibCard(libCard) {
		return fmt.Errorf("[!] ERROR! Reader does not have a valid library card")
	}

	overdueBooks, err := rs.ReservationRepo.GetOverdueByReaderID(ctx, reservation.ReaderID)
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

	existingBook, err := rs.BookRepo.GetByID(ctx, reservation.BookID)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error checking book existence: %v", err)
	}

	if existingBook.Rarity != BookRarityCommon {
		return fmt.Errorf("[!] ERROR! This book is not renewed")
	}

	reservation.IssueDate = time.Now()
	reservation.ReturnDate = time.Now().AddDate(0, 0, ReservationExtensionPeriodDays)
	reservation.State = ReservationExtended

	err = rs.ReservationRepo.Update(ctx, reservation)
	if err != nil {
		return fmt.Errorf("[!] ERROR! Error creating reservation: %v", err)
	}

	return nil
}

func (rs *ReaderService) isValidLibCard(libCard *models.LibCardModel) bool {
	if !libCard.ActionStatus {
		return false
	}

	expiryDate := libCard.IssueDate.AddDate(0, 0, libCard.Validity)

	return time.Now().Before(expiryDate)
}
