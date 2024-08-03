package requesters

import (
	"BookSmart/internal/models"
	"BookSmart/internal/ui/cli/handlers"
	"BookSmart/internal/ui/cli/input"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
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
			fmt.Println(err)
			continue
		}

		switch menuItem {
		case 1:
			err = r.ViewReservations(&reservationsID, tokens)
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			err = r.UpdateReservation(&reservationsID, tokens)
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

func (r *Requester) ViewReservations(reservationsID *[]uuid.UUID, tokens *handlers.TokenResponse) error {
	url := "http://localhost:8000/api/reservations"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
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

	var response []*models.ReservationModel
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	printReservations(response)

	for _, reservation := range response {
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

	reservationIDJSON, err := json.Marshal(reservationID)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:8000/api/reservations/%s", reservationID.String())
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reservationIDJSON))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
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

	// Кодирование тела запроса в JSON
	jsonData, err := json.Marshal(bookID)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://localhost:8000/api/reservations"
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

	fmt.Printf("\n\nBook successfully reserved!\n")

	return nil
}
