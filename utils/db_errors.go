package utils

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tofu345/BMGMT/constants"
)

func PrettyDbError(err error, obj_name ...string) string {
	switch e := err.(type) {
	case *pgconn.PgError:
		return e.Message
	}

	if len(obj_name) > 0 {
		return obj_name[0] + " not found"
	}

	if err.Error() == constants.DBNotFound {
		return constants.NotFound
	}

	return err.Error()
}
