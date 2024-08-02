package tmp

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/services/intfServices"
	"BookSmart/internal/ui/cli/input"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const readerMenu = `Reader's menu:
	1 -- view book information
	2 -- find book
	3 -- add book to favorites (now?)
	4 -- reserve book
	5 -- renew book
	6 -- issue library card
	7 -- renew library card
	0 -- log out
`

type ReaderHandler struct {
	readerService      intfServices.IReaderService
	bookService        intfServices.IBookService
	libCardService     intfServices.ILibCardService
	reservationService intfServices.IReservationService
	logger             *logrus.Entry
}

func NewReaderHandler(
	readerService intfServices.IReaderService,
	bookService intfServices.IBookService,
	libCardService intfServices.ILibCardService,
	reservationService intfServices.IReservationService,
	logger *logrus.Entry,
) *ReaderHandler {
	return &ReaderHandler{
		readerService:      readerService,
		bookService:        bookService,
		libCardService:     libCardService,
		reservationService: reservationService,
		logger:             logger,
	}
}

func (rh *ReaderHandler) Create() error {
	fio, err := input.Fio()
	if err != nil {
		return err
	}

	phoneNumber, err := input.PhoneNumber()
	if err != nil {
		return err
	}

	age, err := input.Age()
	if err != nil {
		return err
	}

	password, err := input.Password()
	if err != nil {
		return err
	}

	reader := &models.ReaderModel{
		ID:          uuid.New(),
		Fio:         fio,
		PhoneNumber: phoneNumber,
		Age:         age,
		Password:    password,
	}

	err = rh.readerService.SignUp(context.Background(), reader)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nRegistration completed successfully!\n")

	return nil
}

func (rh *ReaderHandler) SignIn() error {
	phoneNumber, err := input.PhoneNumber()
	if err != nil {
		return err
	}
	password, err := input.Password()
	if err != nil {
		return err
	}

	readerDTO := &dto.ReaderSignInDTO{
		PhoneNumber: phoneNumber,
		Password:    password,
	}

	_, err = rh.readerService.SignIn(context.Background(), readerDTO)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nAuthentication successful!\n")

	err = rh.readerRequestsHandler()
	if err != nil {
		return err
	}

	return nil
}

func (rh *ReaderHandler) readerRequestsHandler() error {
	for {
		fmt.Printf("\n\n%s", readerMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch menuItem {
		case 0:
			return nil
		default:
			fmt.Printf("\nНеверный пункт меню!\n\n")
		}
	}
}
