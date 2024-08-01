package cli

import (
	"BookSmart/internal/ui/cli/handlers"
	"BookSmart/internal/ui/cli/input"
	"fmt"
	"os"
)

const mainMenu = `Main menu:
	1 -- sign up
	2 -- sign in as reader
	3 -- sign in as administrator
	4 -- view book information
	5 -- find a book
	0 -- exit program
`

type Server struct {
	bookHandler        *handlers.BookHandler
	libCardHandler     *handlers.LibCardHandler
	readerHandler      *handlers.ReaderHandler
	reservationHandler *handlers.ReservationHandler
}

func NewServer(
	bookHandler *handlers.BookHandler,
	libCardHandler *handlers.LibCardHandler,
	readerHandler *handlers.ReaderHandler,
	reservationHandler *handlers.ReservationHandler,
) *Server {
	return &Server{
		bookHandler:        bookHandler,
		libCardHandler:     libCardHandler,
		readerHandler:      readerHandler,
		reservationHandler: reservationHandler,
	}
}

func (s *Server) Run() {
	for {
		fmt.Printf("\n\n%s", mainMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch menuItem {
		case 1:
			err = s.readerHandler.Create()
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			err = s.readerHandler.SignIn()
			if err != nil {
				fmt.Println(err)
			}
		case 0:
			os.Exit(0)
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}
