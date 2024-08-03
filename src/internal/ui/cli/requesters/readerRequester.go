package requesters

import (
	"BookSmart/internal/dto"
	"BookSmart/internal/models"
	"BookSmart/internal/ui/cli/handlers"
	"BookSmart/internal/ui/cli/input"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

func (r *Requester) SignUp() error {
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

	url := "http://localhost:8000/auth/sign-up"
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

func (r *Requester) SignInAsReader(tokens *handlers.TokenResponse) error {
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

	readerDTOJSON, err := json.Marshal(readerDTO)
	if err != nil {
		return err
	}

	url := "http://localhost:8000/auth/sign-in"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(readerDTOJSON))
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

	err = json.NewDecoder(resp.Body).Decode(tokens)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nAuthentication successful!\n")

	return nil
}
