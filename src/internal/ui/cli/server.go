package cli

import (
	"BookSmart/internal/ui/cli/handlers"
	"fmt"
)

const mainMenu = `Меню гостя:
	1 -- зарегистрироваться
	2 -- войти как читатель
	3 -- войти как администратор
	4 -- посмотреть информацию о книге
	5 -- найти книгу
	0 -- выйти
Выберите пункт меню: `

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

	var menuItem int

	for {
		fmt.Printf("\n\n%s", mainMenu)

		_, err := fmt.Scanf("%d", &menuItem)
		if err != nil {
			fmt.Printf("\nПункт меню введён некорректно!\n\n")
			continue
		}
	}
}
