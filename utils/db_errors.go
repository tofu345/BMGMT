package utils

import "github.com/jackc/pgx/v5/pgconn"

func PrettyDbError(err error) string {
	switch e := err.(type) {
	case *pgconn.PgError:
		return e.Message
	}

	return err.Error()
}
