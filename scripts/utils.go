package scripts

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/tofu345/BMGMT/utils"
	"golang.org/x/term"
)

func getScript(name string) (Script, error) {
	for _, v := range scripts {
		if v.name == name {
			return v, nil
		}
	}

	return Script{}, fmt.Errorf("Script '%v' not found", name)
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	text, _ := r.ReadString('\n')
	return strings.TrimSpace(text)
}

func readPassword() string {
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	return string(bytePassword)
}

func getAndComparePasswords() string {
	fmt.Print("> Password: ")
	password := readPassword()

	fmt.Print("> Retype Password: ")
	password2 := readPassword()

	if password != password2 {
		fmt.Println("! Passwords do not match")
		return getAndComparePasswords()
	}

	return password
}

func printValidationErrors(errs utils.ValidationErrs) {
	longest := 0
	for k := range errs.Errors {
		if len(k) > longest {
			longest = len(k)
		}
	}

	for k, v := range errs.Errors {
		fmt.Printf("! %v", k)
		if len(k) < longest {
			fmt.Printf(strings.Repeat(" ", longest-len(k)))
		}
		fmt.Printf("\t%v\n", v)
	}
}
