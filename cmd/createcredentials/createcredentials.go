package main

import (
	"fmt"
	"os"

	"github.com/mailslurper/mailslurper/pkg/auth/basicauth"
)

func main() {
	var userName string
	var password string
	passwordService := basicauth.PasswordService{}

	var hashedPassword []byte
	var err error

	fmt.Printf("Please enter a user name: ")
	if _, err = fmt.Scan(&userName); err != nil {
		fmt.Printf("\nThere was an error: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Printf("Now enter this user's password: ")
	if _, err = fmt.Scan(&password); err != nil {
		fmt.Printf("\nThere was an error: %s\n", err.Error())
		os.Exit(-1)
	}

	if hashedPassword, err = passwordService.HashPassword([]byte(password)); err != nil {
		fmt.Printf("\nThere was an error: %s\n", err.Error())
		os.Exit(-1)
	}

	fmt.Printf("\n\nUser: %s\nPassword: %s\n\n", userName, string(hashedPassword))
}
