package requesters

import (
	"BookSmart-services/core/dto"
	"BookSmart-techUI/input"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const mainMenu = `Main menu:
	1 -- sign up
	2 -- sign in as reader
	3 -- sign in as administrator
	4 -- view books catalog
	0 -- exit program
`

type Requester struct {
	logger          *logrus.Entry
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewRequester(
	logger *logrus.Entry,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *Requester {
	return &Requester{
		logger:          logger,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (r *Requester) Run() {
	for {
		fmt.Printf("\n\n%s", mainMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Printf("\n\n%s\n", err.Error())
			continue
		}

		switch menuItem {
		case 1:
			if err = r.SignUp(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 2:
			if err = r.ProcessReaderActions(); err != nil {
				continue
			}
		case 3:
			if err = r.ProcessAdminActions(); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 4:
			if err = r.ProcessBookCatalogActions(&dto.ReaderTokensDTO{}); err != nil {
				fmt.Printf("\n\n%s\n", err.Error())
			}
		case 0:
			os.Exit(0)
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}
