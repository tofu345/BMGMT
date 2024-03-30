package main

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/tofu345/BMGMT/db"
	"github.com/tofu345/BMGMT/scripts"
	"github.com/tofu345/BMGMT/sqlc"
	"github.com/tofu345/BMGMT/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db.Conn, err = pgx.Connect(db.Ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Conn.Close(db.Ctx)
	db.Q = sqlc.New(db.Conn)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "shell":
			if len(os.Args) >= 2 {
				scripts.Shell(os.Args[2:]...)
			} else {
				scripts.Shell()
			}
		default:
			log.Fatalf("Unknown verb: %v", os.Args[1])
		}
		return
	}

	e := echo.New()
	e.Validator = &utils.Validator
	e.Use(utils.JwtMiddleware)

	registerRoutes(e)

	e.HideBanner = true
	e.Logger.SetLevel(echoLog.INFO)
	e.Logger.Fatal(e.Start(":8000"))
}
