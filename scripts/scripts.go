package scripts

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/tofu345/BMGMT/db"
	"github.com/tofu345/BMGMT/handlers"
	"github.com/tofu345/BMGMT/sqlc"
	"github.com/tofu345/BMGMT/utils"
)

type Script struct {
	name        string
	description string
	function    func()
}

var (
	loggedInAdmin sqlc.User

	r       = bufio.NewReader(os.Stdin)
	scripts = []Script{
		{"create_admin", "Create admin user", createAdmin},
	}
)

func Shell(args ...string) {
	if len(args) > 0 {
		script, err := getScript(args[0])
		if err != nil {
			log.Fatal(err)
		}
		script.function()
		return
	}

	fmt.Println("? 'list' to view commands")
	for {
		input := getUserInput("> ")
		switch input {
		case "":
			continue
		case "help":
			fmt.Println("list\tlist all commands")
			fmt.Println("exit\tquit")
		case "list":
			if len(scripts) == 0 {
				fmt.Println("! There are no scripts")
				return
			}

			for _, script := range scripts {
				fmt.Printf("%v\t%v\n", script.name, script.description)
			}
		case "ex", "exit":
			return
		default:
			script, err := getScript(input)
			if err != nil {
				fmt.Printf("! %v\n", err)
				continue
			}

			script.function()
		}
	}
}

func createAdmin() {
	fmt.Println(">> Create admin user")

	first_name := getUserInput("> First Name: ")
	last_name := getUserInput("> Last Name: ")
	email := getUserInput("> Email: ")
	password := getAndComparePasswords()

	user := handlers.UserData{
		FirstName: first_name,
		LastName:  last_name,
		Password:  password,
		Email:     email,
	}

	errs := utils.Validator.Validate(user)
	if errs != nil {
		printValidationErrors(errs.(utils.ValidationErrs))
		return
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Q.CreateUser(db.Ctx, sqlc.CreateUserParams{
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Password:    hash,
		IsSuperuser: pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("! Admin %v created\n", user.Email)
}
