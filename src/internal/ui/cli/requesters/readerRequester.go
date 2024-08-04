package requesters

import (
	"BookSmart-ui/cli/handlers"
	"BookSmart-ui/cli/input"
	"encoding/json"
	"errors"
	"fmt"
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
	readerSignUpDTO, err := input.SignUpParams()
	if err != nil {
		return err
	}

	request := HTTPRequest{
		Method: "POST",
		URL:    "http://localhost:8000/auth/sign-up",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    readerSignUpDTO,
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		return errors.New(response.Status)
	}

	fmt.Printf("\n\nRegistration completed successfully!\n")

	return nil
}

func (r *Requester) SignIn(tokens *handlers.TokenResponse, stopRefresh <-chan struct{}) error {
	readerSignInDTO, err := input.SignInParams()
	if err != nil {
		return err
	}

	request := HTTPRequest{
		Method: "POST",
		URL:    "http://localhost:8000/auth/sign-in",
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
		return errors.New(response.Status)
	}

	err = json.Unmarshal(response.Body, tokens)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nAuthentication successful!\n")

	go r.Refreshing(tokens, r.accessTokenTTL-time.Second, stopRefresh)

	return nil
}

func (r *Requester) Refresh(tokens *handlers.TokenResponse) error {
	request := HTTPRequest{
		Method: "POST",
		URL:    "http://localhost:8000/auth/refresh",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    tokens.RefreshToken,
		Timeout: 10 * time.Second,
	}

	response, err := SendRequest(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}

	err = json.Unmarshal(response.Body, tokens)
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
