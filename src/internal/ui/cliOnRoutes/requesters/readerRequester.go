package requesters

import (
	"BookSmart/internal/models"
	"BookSmart/internal/ui/cli/input"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
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

type ReaderRequester struct {
	logger *logrus.Entry
}

func NewReaderRequester(
	logger *logrus.Entry,
) *ReaderRequester {
	return &ReaderRequester{
		logger: logger,
	}
}

func (rh *ReaderRequester) Create() error {
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

	readerJSON, err := json.Marshal(reader)
	if err != nil {
		return err
	}

	url := "http://localhost:8000/api/readers/sign-up"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(readerJSON))
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

	// print the response
	var response string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New(response)
	}

	fmt.Printf("\n\nRegistration completed successfully!\n")

	return nil
}

func (rh *ReaderRequester) SignIn() error {
	//phoneNumber, err := input.PhoneNumber()
	//if err != nil {
	//	return err
	//}
	//password, err := input.Password()
	//if err != nil {
	//	return err
	//}

	//readerDTO := &dto.ReaderSignInDTO{
	//	PhoneNumber: phoneNumber,
	//	Password:    password,
	//}

	return nil
}

func (rh *ReaderRequester) readerRequestsHandler() error {
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
