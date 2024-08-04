package requesters

import (
	"BookSmart-services/models"
	"BookSmart-ui/cli/handlers"
	"BookSmart-ui/cli/input"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const libCardMenu = `Library card menu:
	1 -- create library card
	2 -- update library card
	3 -- view info library card
	0 -- go to main menu
`

func (r *Requester) ProcessLibCardActions(tokens *handlers.TokenResponse) error {

	for {
		fmt.Printf("\n\n%s", libCardMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch menuItem {
		case 1:
			err = r.CreateLibCard(tokens)
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			err = r.UpdateLibCard(tokens)
			if err != nil {
				fmt.Println(err)
			}
		case 3:
			err = r.ViewLibCard(tokens)
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

func (r *Requester) CreateLibCard(tokens *handlers.TokenResponse) error {
	request := HTTPRequest{
		Method: "POST",
		URL:    "http://localhost:8000/api/lib-cards",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
		},
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

	fmt.Printf("\n\nSuccessfully created library card!\n")

	return nil
}

func (r *Requester) UpdateLibCard(tokens *handlers.TokenResponse) error {
	request := HTTPRequest{
		Method: "PUT",
		URL:    "http://localhost:8000/api/lib-cards",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
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

	fmt.Printf("\n\nSuccessfully updated library card!\n")

	return nil
}

func (r *Requester) ViewLibCard(tokens *handlers.TokenResponse) error {
	request := HTTPRequest{
		Method: "GET",
		URL:    "http://localhost:8000/api/lib-cards",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
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

	var libCard *models.LibCardModel
	err = json.Unmarshal(response.Body, &libCard)
	if err != nil {
		log.Fatal(err)
	}

	printLibCard(libCard)

	return nil

}

func printLibCard(libCard *models.LibCardModel) {
	issueDateStr := libCard.IssueDate.Format("02.01.2006")

	statusStr := "Inactive"
	if libCard.ActionStatus {
		statusStr = "Active"
	}

	fmt.Printf("\n\nLibrary card:\n")
	fmt.Println(strings.Repeat("-", 27))
	fmt.Printf("Number:     %s\n", libCard.LibCardNum)
	fmt.Printf("Validity:   %d\n", libCard.Validity)
	fmt.Printf("Issue date: %s\n", issueDateStr)
	fmt.Printf("Status:     %s\n", statusStr)
}
