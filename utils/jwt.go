package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/tofu345/BMGMT/constants"
	"github.com/tofu345/BMGMT/db"
	"github.com/tofu345/BMGMT/sqlc"
)

func LoginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("user") == nil {
			return c.String(http.StatusUnauthorized, constants.MissingPerms)
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
			return c.String(http.StatusBadRequest, err.Error())
		}
		c.Set("user", user)
		return next(c)
	}
}

func jwtAuth(token string) (sqlc.User, error) {
	if token == "" {
		return sqlc.User{}, constants.ErrMissingToken
	}

	tokens := strings.Split(token, " ")
	if len(tokens) != 2 {
		return sqlc.User{}, constants.ErrInvalidToken
	}

	payload, err := DecodeToken(tokens[1])
	if err != nil {
		return sqlc.User{}, constants.ErrInvalidToken
	}

	email := payload["email"]
	switch email := email.(type) {
	case string:
		user, err := db.Q.GetUserByEmail(db.Ctx, email)
		if err != nil {
			return sqlc.User{}, err
		}

		return user, nil
	default:
		return sqlc.User{}, constants.ErrInvalidToken
	}
}

func defaultJwtClaims(user sqlc.User) jwt.MapClaims {
	time_now := time.Now()
	return jwt.MapClaims{
		"iss":   constants.JWT_ISSUER,
		"iat":   time_now.Unix(),
		"exp":   time_now.Add(time.Hour).Unix(),
		"email": user.Email,
	}
}

func AccessToken(user sqlc.User) (string, error) {
	key := []byte(constants.JWT_KEY)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, defaultJwtClaims(user))
	return t.SignedString(key)
}

func RefreshToken(user sqlc.User) (string, error) {
	key := []byte(constants.JWT_KEY)
	claims := defaultJwtClaims(user)
	claims["ref"] = true
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(key)
}

func DecodeToken(tokenData string) (map[string]any, error) {
	token, err := jwt.Parse(tokenData, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(constants.JWT_KEY), nil
	})
	if err != nil {
		return map[string]any{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return map[string]any{}, err
	}
}
