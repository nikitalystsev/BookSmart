package requesters

import (
	"BookSmart-services/dto"
	"BookSmart-services/models"
	"BookSmart-ui/cli/handlers"
	"BookSmart-ui/cli/input"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strings"
)

/*
	3 -- add book to favorites (now?)
	4 -- reserve book
	5 -- renew book
	6 -- issue library card
	7 -- renew library card
*/

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
		case 0:
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}
func (r *Requester) viewFirstPage(bookPagesID *[]uuid.UUID) (dto.BookParamsDTO, error) {
	isWithParams, err := input.IsWithParams()
	if err != nil {
		return dto.BookParamsDTO{Limit: PageLimit, Offset: 0}, err
	}

	*bookPagesID = make([]uuid.UUID, 0)

	var params dto.BookParamsDTO
	if !isWithParams {
		params, err = input.Params()
		if err != nil {
			return dto.BookParamsDTO{Limit: PageLimit, Offset: 0}, err
		}
		params.Limit = PageLimit
		params.Offset = 0
	} else {
		params = dto.BookParamsDTO{Limit: PageLimit, Offset: 0}
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return dto.BookParamsDTO{Limit: PageLimit, Offset: 0}, err
	}

	url := "http://localhost:8000/general/books"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(paramsJSON))
	if err != nil {
		return dto.BookParamsDTO{Limit: PageLimit, Offset: 0}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return dto.BookParamsDTO{Limit: PageLimit, Offset: 0}, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var response string
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return dto.BookParamsDTO{Limit: PageLimit, Offset: 0}, err
		}
		return dto.BookParamsDTO{Limit: PageLimit, Offset: 0}, errors.New(response)
	}

	var response []*models.BookModel
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	printBooks(response, 0)
	updateParams(&params, bookPagesID, response)

	return params, nil
}

func (r *Requester) viewNextPage(params *dto.BookParamsDTO, bookPagesID *[]uuid.UUID) error {
	if params.Limit == 0 {
		params.Limit = PageLimit
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return err
	}

	url := "http://localhost:8000/general/books"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(paramsJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var response string
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return err
		}
		return errors.New(response)
	}

	var response []*models.BookModel
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	printBooks(response, params.Offset)
	updateParams(params, bookPagesID, response)

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

	url := fmt.Sprintf("http://localhost:8000/general/books/%s", bookID.String())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println("error closing body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var response string
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return err
		}
		return errors.New(response)
	}

	var response *models.BookModel
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	printBook(response, num)

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

	// Кодирование тела запроса в JSON
	jsonData, err := json.Marshal(bookID)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://localhost:8000/api/favorites"
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

	fmt.Printf("\n\nBook successfully added to your favorites!\n")

	return nil
}

// Функция для красивого вывода информации о книге
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

// Функция для вывода элементов слайса указателей на BookModel
func printBooks(response []*models.BookModel, offset int) {
	// Определим ширину для выравнивания
	titleWidth := 60
	authorWidth := 60

	// Выведем заголовок таблицы
	fmt.Printf("%-5s %-60s %-60s\n", "No.", "Title", "Author")
	fmt.Println(strings.Repeat("-", 5+1+titleWidth+1+authorWidth))

	for i, book := range response {
		fmt.Printf("%-5d %-60s %-60s\n", offset+i, truncate(book.Title, titleWidth), truncate(book.Author, authorWidth))
	}
}

// Вспомогательная функция для обрезки строки до заданной длины
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
