package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tofu345/BMGMT/sqlc"
)

var ConnPool *pgxpool.Pool
var Ctx context.Context
var Q *sqlc.Queries

func init() {
	Ctx = context.Background()
}
