package requesters

import (
	"BookSmart-services/core/dto"
	"BookSmart-services/core/models"
	"BookSmart-techUI/handlers"
	"BookSmart-techUI/input"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

const catalogMenu = `Catalog's menu:
	1 -- view books
	2 -- next page
	3 -- view info about book
	4 -- add book to favorites
	5 -- reserve book
	0 -- go to main menu
`

const PageLimit = 10

func (r *Requester) ProcessBookCatalogActions(tokens *handlers.TokenResponse) error {
	var params dto.BookParamsDTO
	var bookPagesID []uuid.UUID // массив id выведенных книг

	for {
		fmt.Printf("\n\n%s", catalogMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Printf("\n\n%s\n", err.Error())
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
		case 0:
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}
func (r *Requester) viewFirstPage(params *dto.BookParamsDTO, bookPagesID *[]uuid.UUID) error {
	isWithParams, err := input.IsWithParams()
	if err != nil {
		return err
	}

	*bookPagesID = make([]uuid.UUID, 0)

	if !isWithParams {
		if *params, err = input.Params(); err != nil {
			return err
		}
		params.Limit = PageLimit
		params.Offset = 0
	} else {
		*params = dto.BookParamsDTO{Limit: PageLimit, Offset: 0}
	}

	request := HTTPRequest{
		Method: "POST",
		URL:    "http://localhost:8000/general/books",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    *params,
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

	var books []*models.BookModel
	if err = json.Unmarshal(response.Body, &books); err != nil {
		return err
	}

	printBooks(books, 0)
	updateParams(params, bookPagesID, books)

	return nil
}

func (r *Requester) viewNextPage(params *dto.BookParamsDTO, bookPagesID *[]uuid.UUID) error {
	if params.Limit == 0 {
		params.Limit = PageLimit
	}

	request := HTTPRequest{
		Method: "POST",
		URL:    "http://localhost:8000/general/books",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    *params,
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

	var books []*models.BookModel
	if err = json.Unmarshal(response.Body, &books); err != nil {
		return err
	}

	printBooks(books, params.Offset)
	updateParams(params, bookPagesID, books)

	return nil
}

func (r *Requester) ViewBook(bookPagesID *[]uuid.UUID) error {
	num, err := input.BookPagesNumber()
	if err != nil {
		return err
	}

	if num > len(*bookPagesID) || num < 0 {
		return errors.New("book number out of range")
	}

	bookID := (*bookPagesID)[num]

	request := HTTPRequest{
		Method: "GET",
		URL:    fmt.Sprintf("http://localhost:8000/general/books/%s", bookID.String()),
		Headers: map[string]string{
			"Content-Type": "application/json",
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

	var book *models.BookModel
	if err = json.Unmarshal(response.Body, &book); err != nil {
		return err
	}

	printBook(book, num)

	return nil

}

func (r *Requester) AddToFavorites(bookPagesID *[]uuid.UUID, accessToken string) error {
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
		URL:    "http://localhost:8000/api/favorites",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
		},
		Body:    bookID,
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return errors.New("you are not authenticated")
	}
	if response.StatusCode != http.StatusCreated {
		var info string
		if err = json.Unmarshal(response.Body, &info); err != nil {
			return err
		}
		return errors.New(info)
	}

	fmt.Printf("\n\nBook successfully added to your favorites!\n")

	return nil
}

func printBook(book *models.BookModel, num int) {
	fmt.Printf("\n\nBook №%d:\n", num)
	fmt.Println(strings.Repeat("-", 150))
	fmt.Printf("Title:           %s\n", book.Title)
	fmt.Printf("Author:          %s\n", book.Author)
	fmt.Printf("Publisher:       %s\n", book.Publisher)
	fmt.Printf("Copies Number:   %d\n", book.CopiesNumber)
	fmt.Printf("Rarity:          %s\n", book.Rarity)
	fmt.Printf("Genre:           %s\n", book.Genre)
	fmt.Printf("Publishing Year: %d\n", book.PublishingYear)
	fmt.Printf("Language:        %s\n", book.Language)
	fmt.Printf("Age Limit:       %d\n", book.AgeLimit)
}

func printBooks(response []*models.BookModel, offset int) {
	titleWidth := 60
	authorWidth := 60

	fmt.Printf("%-5s %-60s %-60s\n", "No.", "Title", "Author")
	fmt.Println(strings.Repeat("-", 5+1+titleWidth+1+authorWidth))

	for i, book := range response {
		fmt.Printf("%-5d %-60s %-60s\n", offset+i, truncate(book.Title, titleWidth), truncate(book.Author, authorWidth))
	}
}

func truncate(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength-3] + "..."
	}
	return s
}

func updateParams(params *dto.BookParamsDTO, bookPagesID *[]uuid.UUID, books []*models.BookModel) {
	params.Offset += PageLimit

	for _, book := range books {
		*bookPagesID = append(*bookPagesID, book.ID)
	}
}
