package requesters

import (
	"BookSmart-services/core/models"
	"BookSmart-ui/handlers"
	"BookSmart-ui/input"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

const reservationsMenu = `Reservations menu:
	1 -- view your reservations 
	2 -- update your reservation
	0 -- go to main menu
`

func (r *Requester) ProcessReservationsActions(tokens *handlers.TokenResponse) error {
	var reservationsID []uuid.UUID // массив id выведенных броней

	for {
		fmt.Printf("\n\n%s", reservationsMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Printf("\n\n%s\n", err.Error())
			continue
		}

		switch menuItem {
		case 1:
			if err = r.ViewReservations(&reservationsID, tokens); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 2:
			if err = r.UpdateReservation(&reservationsID, tokens); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 0:
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}

func (r *Requester) ViewReservations(reservationsID *[]uuid.UUID, tokens *handlers.TokenResponse) error {
	request := HTTPRequest{
		Method: "GET",
		URL:    "http://localhost:8000/api/reservations",
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

	var reservations []*models.ReservationModel
	if err = json.Unmarshal(response.Body, &reservations); err != nil {
		return err
	}

	printReservations(reservations)

	for _, reservation := range reservations {
		*reservationsID = append(*reservationsID, reservation.ID)
	}

	return nil
}

func (r *Requester) UpdateReservation(reservationsID *[]uuid.UUID, tokens *handlers.TokenResponse) error {
	num, err := input.ReservationNumber()
	if err != nil {
		return err
	}

	if num > len(*reservationsID) || num < 0 {
		return errors.New("reservation number out of range")
	}

	reservationID := (*reservationsID)[num]

	request := HTTPRequest{
		Method: "PUT",
		URL:    fmt.Sprintf("http://localhost:8000/api/reservations/%s", reservationID.String()),
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", tokens.AccessToken),
		},
		Body:    reservationID,
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

	fmt.Printf("\n\nReservation successfully updated!\n")

	return nil
}

func printReservations(reservations []*models.ReservationModel) {
	fmt.Printf("\n\n%-5s %-12s %-12s %-10s\n", "No.", "Issue Date", "Return Date", "State")
	fmt.Println(strings.Repeat("-", 40))

	for i, r := range reservations {
		fmt.Printf("%-5d %-12s %-12s %-10s\n", i, r.IssueDate.Format("2006-01-02"), r.ReturnDate.Format("2006-01-02"), r.State)
	}
}

func (r *Requester) ReserveBook(bookPagesID *[]uuid.UUID, accessToken string) error {
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
		URL:    "http://localhost:8000/api/reservations",
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

	fmt.Printf("\n\nBook successfully reserved!\n")

	return nil
}
