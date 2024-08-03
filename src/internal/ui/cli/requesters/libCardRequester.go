package requesters

import (
	"BookSmart/internal/models"
	"BookSmart/internal/ui/cli/handlers"
	"BookSmart/internal/ui/cli/input"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	url := "http://localhost:8000/api/lib-cards"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
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

	if resp.StatusCode != http.StatusCreated {
		var response string
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return err
		}
		return errors.New(response)
	}

	fmt.Printf("\n\nSuccessfully created library card!\n")

	return nil
}

func (r *Requester) UpdateLibCard(tokens *handlers.TokenResponse) error {
	url := "http://localhost:8000/api/lib-cards"
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
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

	fmt.Printf("\n\nSuccessfully updated library card!\n")

	return nil
}

func (r *Requester) ViewLibCard(tokens *handlers.TokenResponse) error {

	url := "http://localhost:8000/api/lib-cards"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
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

	var response *models.LibCardModel
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	printLibCard(response)

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
