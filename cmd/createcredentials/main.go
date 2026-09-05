// Command createcredentials bcrypt-hashes a username/password pair for
// pasting into config.json's "credentials" map.
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/p0vidl0/mylslurper/internal/auth"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter a user name: ")
	userName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nThere was an error: %s\n", err)
		os.Exit(1)
	}

	fmt.Print("Now enter this user's password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\nThere was an error: %s\n", err)
		os.Exit(1)
	}

	hashedPassword, err := auth.HashPassword([]byte(trim(password)))
	if err != nil {
		fmt.Printf("\nThere was an error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n\nUser: %s\nPassword: %s\n\n", trim(userName), string(hashedPassword))
}

func trim(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}
