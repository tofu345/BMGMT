package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/tofu345/BMGMT/sqlc"
)

var Conn *pgx.Conn
var Ctx context.Context
var Q *sqlc.Queries

func init() {
	Ctx = context.Background()
}
