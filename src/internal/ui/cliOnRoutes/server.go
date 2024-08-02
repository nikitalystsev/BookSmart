package cliOnRoutes

import (
	"BookSmart/internal/ui/cliOnRoutes/input"
	"BookSmart/internal/ui/cliOnRoutes/requesters"
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
	bookRequester        *requesters.BookRequester
	libCardRequester     *requesters.LibCardRequester
	readerRequester      *requesters.ReaderRequester
	reservationRequester *requesters.ReservationRequester
}

func NewServer(
	bookRequester *requesters.BookRequester,
	libCardRequester *requesters.LibCardRequester,
	readerRequester *requesters.ReaderRequester,
	reservationRequester *requesters.ReservationRequester,
) *Server {
	return &Server{
		bookRequester:        bookRequester,
		libCardRequester:     libCardRequester,
		readerRequester:      readerRequester,
		reservationRequester: reservationRequester,
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
			err = s.readerRequester.Create()
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			err = s.readerRequester.SignIn()
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
