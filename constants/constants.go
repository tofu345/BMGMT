package constants

import (
	"errors"
	"os"
	"strings"
)

func init() {
	JWT_KEY = os.Getenv("JWT_KEY")
	ALLOWED_HOSTS = strings.Split(os.Getenv("ALLOWED_HOSTS"), ",")
}

var (
	JWT_KEY       string
	ALLOWED_HOSTS []string

	ErrMissingToken = errors.New("missing token")
	ErrInvalidToken = errors.New(InvalidToken)
	ErrInvalidData  = errors.New(InvalidData)
	ErrExpiredToken = errors.New("expired token")
	ErrNotFound     = errors.New(NotFound)
)

const (
	JWT_ISSUER = "BMGMT"

	MissingPerms = "missing permissions"
	Unauthorized = "unauthorized"
	InvalidData  = "invalid data"
	BadRequest   = "bad request"
	InvalidToken = "invalid token"
	NotFound     = "object not found"
	DBNotFound   = "no rows in result set"
)
