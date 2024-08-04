package requesters

import (
	"BookSmart-services/dto"
	"BookSmart-ui/cli/handlers"
	"BookSmart-ui/cli/input"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
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
			params, err = r.viewFirstPage(&bookPagesID)
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
	// Кодирование тела запроса в JSON
	jsonData, err := json.Marshal(newBook)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://localhost:8000/api/admin/books"
	// Создание нового HTTP-запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		var response string
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return err
		}
		return errors.New(response)
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

	url := fmt.Sprintf("http://localhost:8000/api/admin/books/%s", bookID.String())
	// Создание нового HTTP-запроса
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var response string
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return err
		}
		return errors.New(response)
	}

	fmt.Printf("\n\nBook successfully deleted!\n")

	return nil
}
