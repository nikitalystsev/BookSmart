package requesters

import (
	"BookSmart-services/dto"
	"BookSmart-ui/cli/handlers"
	"BookSmart-ui/cli/input"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (r *Requester) ProcessAdminActions() error {
	var tokens handlers.TokenResponse
	stopRefresh := make(chan struct{})

	if err := r.SignIn(&tokens, stopRefresh); err != nil {
		fmt.Println(err)
		return err
	}

	for {
		fmt.Printf("\n\n%s", readerMainMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch menuItem {
		case 1:
			err = r.ProcessAdminBookCatalogActions(&tokens)
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			err = r.ProcessLibCardActions(&tokens)
			if err != nil {
				fmt.Println(err)
			}
		case 3:
			err = r.ProcessReservationsActions(&tokens)
			if err != nil {
				fmt.Println(err)
			}
		case 0:
			close(stopRefresh)
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}

const adminCatalogMenu = `Admin's Catalog menu:
	1 -- view books
	2 -- next page
	3 -- view info about book
	4 -- add book to favorites
	5 -- reserve book
	6 -- add new book
	7 -- delete book
	0 -- go to main menu
`

func (r *Requester) ProcessAdminBookCatalogActions(tokens *handlers.TokenResponse) error {
	var params dto.BookParamsDTO
	var bookPagesID []uuid.UUID // массив id выведенных книг

	for {
		fmt.Printf("\n\n%s", adminCatalogMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch menuItem {
		case 1:
			err = r.viewFirstPage(&params, &bookPagesID)
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			err = r.viewNextPage(&params, &bookPagesID)
			if err != nil {
				fmt.Println(err)
			}
		case 3:
			err = r.ViewBook(&bookPagesID)
			if err != nil {
				fmt.Println(err)
			}
		case 4:
			err = r.AddToFavorites(&bookPagesID, tokens.AccessToken)
			if err != nil {
				fmt.Println(err)
			}
		case 5:
			err = r.ReserveBook(&bookPagesID, tokens.AccessToken)
			if err != nil {
				fmt.Println(err)
			}
		case 6:
			err = r.AddNewBook(tokens.AccessToken)
			if err != nil {
				fmt.Println(err)
			}
		case 7:
			err = r.DeleteBook(&bookPagesID, tokens.AccessToken)
			if err != nil {
				fmt.Println(err)
			}
		case 0:
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}

func (r *Requester) AddNewBook(accessToken string) error {
	newBook, err := input.Book()
	if err != nil {
		return err
	}

	request := HTTPRequest{
		Method: "POST",
		URL:    "http://localhost:8000/api/admin/books",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		},
		Body:    newBook,
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	fmt.Printf("\n\nBook successfully created!\n")

	return nil
}

func (r *Requester) DeleteBook(bookPagesID *[]uuid.UUID, accessToken string) error {
	num, err := input.BookPagesNumber()
	if err != nil {
		return err
	}

	if num > len(*bookPagesID) || num < 0 {
		return errors.New("book number out of range")
	}

	bookID := (*bookPagesID)[num]

	request := HTTPRequest{
		Method: "POST",
		URL:    fmt.Sprintf("http://localhost:8000/api/admin/books/%s", bookID.String()),
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		},
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	fmt.Printf("\n\nBook successfully deleted!\n")

	return nil
}
