package repository

import (
    "context"
    "finances/internal/entity"
    "github.com/jackc/pgx/v4"
)

type TransactionRepository interface {
    CreateTransaction(ctx context.Context, tx pgx.Tx, transaction entity.Transaction) error
    GetLastTransactions(ctx context.Context, userID int64, limit int) ([]entity.Transaction, error)
}