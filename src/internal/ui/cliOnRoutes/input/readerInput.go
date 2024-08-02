package input

import (
	"bufio"
	"fmt"
	"github.com/howeyc/gopass"
	"os"
	"strconv"
	"strings"
)

func Fio() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input your FIO: ")

	fio, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	fio = strings.TrimSpace(fio)

	return fio, nil
}

func PhoneNumber() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Input your phone number: ")

	phoneNumber, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	phoneNumber = strings.TrimSpace(phoneNumber)

	return phoneNumber, nil
}

func Age() (uint, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input your age: ")

	ageStr, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	ageStr = strings.TrimSpace(ageStr)

	ageInt, err := strconv.Atoi(ageStr)
	if err != nil {
		return 0, err
	}

	age := uint(ageInt)

	return age, nil
}

func Password() (string, error) {
	fmt.Print("Input your password: ")

	silentPassword, err := gopass.GetPasswdMasked()
	if err != nil {
		return "", err
	}

	password := string(silentPassword)
	password = strings.TrimSpace(password)

	return password, nil
}

func MenuItem() (int, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input menu item: ")

	menuItemStr, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	menuItemStr = strings.TrimSpace(menuItemStr)

	menuItemInt, err := strconv.Atoi(menuItemStr)
	if err != nil {
		return 0, err
	}

	return menuItemInt, nil
}
