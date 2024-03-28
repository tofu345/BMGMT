package main

import (
	"log"
	"os"

	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	"github.com/tofu345/BMGMT/db"
	"github.com/tofu345/BMGMT/sqlc"
	"github.com/tofu345/BMGMT/utils"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return utils.FmtValidationErrs(err)
	}
	return nil
}

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

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(utils.JwtMiddleware)

	registerRoutes(e)

	e.HideBanner = true
	e.Logger.SetLevel(echoLog.INFO)
	e.Logger.Fatal(e.Start(":8000"))
}
