package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tofu345/BMGMT/constants"
	"github.com/tofu345/BMGMT/db"
	"github.com/tofu345/BMGMT/sqlc"
	"github.com/tofu345/BMGMT/utils"
)

func GetUsers(c echo.Context) error {
	users, err := db.Q.ListUsers(db.Ctx)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return c.JSONBlob(http.StatusOK, []byte("[]"))
	}
	return c.JSON(http.StatusOK, users)
}

type UserDisplay struct {
	Email     string         `json:"email"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Location  []sqlc.Location `json:"admin_locations,omitempty"`
}

func GetUserInfo(c echo.Context) error {
	user := c.Get("user").(sqlc.User)
	locations, err := db.Q.GetUserLocAdmins(db.Ctx, user.ID)
	if err != nil {
		if err.Error() != constants.DBNotFound {
			return c.String(http.StatusBadRequest, utils.PrettyDbError(err))
		}
	}

	return c.JSON(http.StatusOK, UserDisplay{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Location:  locations,
	})
}

type UserDTO struct {
	Email     string `json:"email" validate:"required,email,max=30"`
	FirstName string `json:"first_name" validate:"required,min=3,max=20"`
	LastName  string `json:"last_name" validate:"required,min=3,max=20"`
	Password  string `json:"password" validate:"required,min=6,max=20"`
}

func CreateUser(c echo.Context) error {
	user := new(UserDTO)
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
		IsSuperuser: false,
	})
	if err != nil {
		return c.String(http.StatusBadRequest, utils.PrettyDbError(err))
	}

	return c.JSON(http.StatusCreated, UserDisplay{
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
	})
}

type GenTokenDTO struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

func GenerateTokenPair(c echo.Context) error {
	data := new(GenTokenDTO)
	if err := c.Bind(data); err != nil {
		return c.String(http.StatusBadRequest, constants.BadRequest)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user, err := db.Q.GetUserByEmail(db.Ctx, data.Email)
	if err != nil {
		return c.String(http.StatusBadRequest, utils.PrettyDbError(err))
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

type RegenTokenDTO struct {
	RefreshToken string `json:"refresh" validate:"required"`
}

func RegenerateAccessToken(c echo.Context) error {
	data := new(RegenTokenDTO)
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
