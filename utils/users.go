package utils

import (
	"github.com/tofu345/BMGMT/constants"
	"github.com/tofu345/BMGMT/db"
	"github.com/tofu345/BMGMT/sqlc"
)

func GetUserByEmail(email string) (sqlc.User, error) {
	users, err := db.Q.GetUserByEmail(db.Ctx, email)
	if err != nil {
		return sqlc.User{}, err
	}
	if len(users) == 0 {
		return sqlc.User{}, constants.ErrNotFound
	}
	return users[0], nil
}
