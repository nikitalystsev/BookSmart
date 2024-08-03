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
	"time"
)

const readerMainMenu = `Main menu:
	1 -- go to books catalog 
	2 -- go to library card
	3 -- go to your reservations
	0 -- log out
`

func (r *Requester) ProcessReaderActions() error {
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
			err = r.ProcessBookCatalogActions(&tokens)
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

func (r *Requester) SignIn(tokens *handlers.TokenResponse, stopRefresh <-chan struct{}) error {
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

	go r.Refreshing(tokens, r.accessTokenTTL-time.Second, stopRefresh)

	return nil
}

func (r *Requester) Refresh(tokens *handlers.TokenResponse) error {
	url := "http://localhost:8000/auth/refresh"

	refreshTokenJSON, err := json.Marshal(tokens.RefreshToken)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(refreshTokenJSON))
	if err != nil {
		return err
	}
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

	err = json.NewDecoder(resp.Body).Decode(tokens)
	if err != nil {
		return err
	}

	//fmt.Printf("\n\nSuccessful refresh tokens!\n")

	return nil
}

func (r *Requester) Refreshing(tokens *handlers.TokenResponse, interval time.Duration, stopRefresh <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := r.Refresh(tokens)
			if err != nil {
				fmt.Printf("Error refreshing tokens: %v\n", err)
			}
		case <-stopRefresh:
			return
		}
	}
}
