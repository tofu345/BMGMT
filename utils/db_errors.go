package utils

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tofu345/BMGMT/constants"
)

func PrettyDbError(err error) string {
	switch e := err.(type) {
	case *pgconn.PgError:
		return e.Message
	}

	if err.Error() == constants.DBNotFound {
		return constants.NotFound
	}
	return err.Error()
}
