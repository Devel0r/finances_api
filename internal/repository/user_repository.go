package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type UserRepository interface {
	GetBalance(ctx context.Context, tx pgx.Tx, userID int64) (float64, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, userID int64, amount float64) error
}
