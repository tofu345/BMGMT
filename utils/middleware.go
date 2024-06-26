package utils

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tofu345/BMGMT/constants"
	"github.com/tofu345/BMGMT/db"
	"github.com/tofu345/BMGMT/sqlc"
)

func LocationAdminRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(sqlc.User)
		if !ok {
			return c.String(http.StatusUnauthorized, constants.Unauthorized)
		}

		if !user.IsSuperuser {
			id, err := strconv.Atoi(c.Param("loc_id"))
			if err != nil {
				return c.String(http.StatusBadRequest, constants.BadRequest)
			}
			loc_id := int64(id)
			admins, err := db.Q.GetLocationAdmins(db.Ctx, loc_id)
			if err != nil {
				return c.String(http.StatusBadRequest, PrettyDbError(err))
			}
			if !StringSliceContains(admins, user.Email) {
				return c.String(http.StatusUnauthorized, constants.Unauthorized)
			}
		}

		return next(c)
	}
}

func SuperUserRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		if user == nil {
			return c.String(http.StatusUnauthorized, constants.Unauthorized)
		}

		if !user.(sqlc.User).IsSuperuser {
			return c.String(http.StatusUnauthorized, constants.MissingPerms)
		}
		return next(c)
	}
}

func LoginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("user") == nil {
			return c.String(http.StatusUnauthorized, constants.Unauthorized)
		}

		return next(c)
	}
}

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return next(c)
		}

		user, err := jwtAuth(token)
		if err != nil {
			log.Println(err.Error())
			return next(c)
		}
		c.Set("user", user)
		return next(c)
	}
}
