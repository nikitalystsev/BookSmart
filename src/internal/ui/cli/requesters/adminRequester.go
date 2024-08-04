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

	if err := r.SignInAsAdmin(&tokens, stopRefresh); err != nil {
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
			fmt.Println("\n\nyou have successfully log out")
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}

func (r *Requester) SignInAsAdmin(tokens *handlers.TokenResponse, stopRefresh <-chan struct{}) error {
	readerSignInDTO, err := input.SignInParams()
	if err != nil {
		return err
	}

	request := HTTPRequest{
		Method: "POST",
		URL:    "http://localhost:8000/auth/admin/sign-in",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    readerSignInDTO,
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

	if err = json.Unmarshal(response.Body, tokens); err != nil {
		return err
	}

	fmt.Printf("\n\nAuthentication successful!\n")

	go r.Refreshing(tokens, r.accessTokenTTL, stopRefresh)

	return nil
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
			if err = r.viewFirstPage(&params, &bookPagesID); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 2:
			if err = r.viewNextPage(&params, &bookPagesID); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 3:
			if err = r.ViewBook(&bookPagesID); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 4:
			if err = r.AddToFavorites(&bookPagesID, tokens.AccessToken); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 5:
			if err = r.ReserveBook(&bookPagesID, tokens.AccessToken); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 6:
			if err = r.AddNewBook(tokens.AccessToken); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 7:
			if err = r.DeleteBook(&bookPagesID, tokens.AccessToken); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
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
