package requesters

import (
	"BookSmart/internal/ui/cli/handlers"
	"BookSmart/internal/ui/cli/input"
	"fmt"
)

const adminMainMenu = `Main menu:
	1 -- go to books catalog 
	2 -- go to library card 
	0 -- log out
`

func (r *Requester) ProcessAdminActions() error {
	var tokens handlers.TokenResponse
	stopRefresh := make(chan struct{})

	if err := r.SignInAsReader(&tokens, stopRefresh); err != nil {
		fmt.Println(err)
		return err
	}

	for {
		fmt.Printf("\n\n%s", adminMainMenu)

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
		case 0:
			close(stopRefresh)
			return nil
		default:
			fmt.Printf("\n\nWrong menu item!\n")
		}
	}
}
