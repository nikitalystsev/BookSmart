package requesters

import (
	"BookSmart/internal/ui/cliOnRoutes/handlers"
	"BookSmart/internal/ui/cliOnRoutes/input"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

const mainMenu = `Main menu:
	1 -- sign up
	2 -- sign in as reader
	3 -- sign in as administrator
	4 -- view books catalog
	0 -- exit program
`

type Requester struct {
	logger *logrus.Entry
}

func NewRequester(logger *logrus.Entry) *Requester {
	return &Requester{logger: logger}
}

func (r *Requester) Run() {
	var tokens handlers.TokenResponse

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
			err = r.SignInAsReader(&tokens)
			if err != nil {
				fmt.Println(err)
			}
		case 4:
			err = r.ProcessBookCatalogActions(&tokens)
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
