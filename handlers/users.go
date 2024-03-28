package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/tofu345/BMGMT/constants"
	"github.com/tofu345/BMGMT/db"
	"github.com/tofu345/BMGMT/sqlc"
	"github.com/tofu345/BMGMT/utils"
)

type UserDisplay struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	IsSuperuser bool   `json:"is_superuser"`
}

func GetUsers(c echo.Context) error {
	user := c.Get("user")
	if !user.(sqlc.User).IsSuperuser.Bool {
		return c.String(http.StatusUnauthorized, constants.MissingPerms)
	}

	users, err := db.Q.ListUsers(db.Ctx)
	if err != nil {
		return err
	}

	out := make([]UserDisplay, len(users))
	for i, u := range users {
		out[i] = UserDisplay{
			Email:       u.Email,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			IsSuperuser: u.IsSuperuser.Bool,
		}
	}

	return c.JSON(http.StatusOK, out)
}

func GetUserInfo(c echo.Context) error {
	user := c.Get("user").(sqlc.User)
	return c.JSON(http.StatusOK, UserDisplay{
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		IsSuperuser: user.IsSuperuser.Bool,
	})
}

type UserData struct {
	Email     string `json:"email" validate:"required,email,max=30"`
	FirstName string `json:"first_name" validate:"required,min=3,max=20"`
	LastName  string `json:"last_name" validate:"required,min=3,max=20"`
	Password  string `json:"password" validate:"required,min=6,max=20"`
}

func CreateUser(c echo.Context) error {
	user := new(UserData)
	if err := c.Bind(user); err != nil {
		return c.String(http.StatusBadRequest, constants.InvalidData)
	}
	if err := c.Validate(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.String(http.StatusBadRequest, "error hashing password")
	}

	newUser, err := db.Q.CreateUser(db.Ctx, sqlc.CreateUserParams{
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Password:    hash,
		IsSuperuser: pgtype.Bool{Bool: false, Valid: true},
	})
	if err != nil {
		return c.String(http.StatusBadRequest, utils.PrettyDbError(err))
	}

	return c.JSON(http.StatusCreated, UserDisplay{
		Email:       newUser.Email,
		FirstName:   newUser.FirstName,
		LastName:    newUser.LastName,
		IsSuperuser: newUser.IsSuperuser.Bool,
	})
}

type GenTokenData struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

func GenerateTokenPair(c echo.Context) error {
	data := new(GenTokenData)
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, constants.BadRequest)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user, err := db.Q.GetUserByEmail(db.Ctx, data.Email)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if !utils.CheckPasswordHash(data.Password, user.Password) {
		return c.String(http.StatusBadRequest, "incorrect password")
	}

	access, err := utils.AccessToken(user)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	refresh, err := utils.RefreshToken(user)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"access": access, "refresh": refresh})
}

type RegenTokenData struct {
	RefreshToken string `json:"refresh" validate:"required"`
}

func RegenerateAccessToken(c echo.Context) error {
	data := new(RegenTokenData)
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, constants.BadRequest)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	payload, err := utils.DecodeToken(data.RefreshToken)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if _, exists := payload["ref"]; !exists {
		return c.String(http.StatusBadRequest, constants.InvalidToken)
	}

	email := payload["email"]
	switch email := email.(type) {
	case string:
		user, err := db.Q.GetUserByEmail(db.Ctx, email)
		if err != nil {
			return c.String(http.StatusBadRequest, utils.PrettyDbError(err))
		}

		access, err := utils.AccessToken(user)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.JSON(http.StatusBadRequest, map[string]string{"access": access})
	}

	return c.String(http.StatusBadRequest, constants.InvalidToken)
}
