package requesters

import (
	"BookSmart-ui/cli/handlers"
	"BookSmart-ui/cli/input"
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

func NewRequester(logger *logrus.Entry, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *Requester {
	return &Requester{logger: logger, accessTokenTTL: accessTokenTTL, refreshTokenTTL: refreshTokenTTL}
}

func (r *Requester) Run() {
	for {
		fmt.Printf("\n\n%s", mainMenu)

		menuItem, err := input.MenuItem()
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch menuItem {
		case 1:
			err = r.SignUp()
			if err != nil {
				fmt.Println(err)
			}
		case 2:
			err = r.ProcessReaderActions()
			if err != nil {
				fmt.Println(err)
			}
		case 3:
			err = r.ProcessAdminActions()
			if err != nil {
				fmt.Println(err)
			}
		case 4:
			err = r.ProcessBookCatalogActions(&handlers.TokenResponse{})
			if err != nil {
				fmt.Println(err)
			}
		case 0:
			os.Exit(0)
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}
