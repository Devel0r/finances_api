package repository

import (
	"context"

	"finances/internal/entity"

	"github.com/jackc/pgx/v4"
)

type transactionRepositoryImpl struct {
	db *pgx.Conn
}

func NewTransactionRepository(db *pgx.Conn) TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) CreateTransaction(ctx context.Context, tx pgx.Tx, transaction entity.Transaction) error {
	_, err := tx.Exec(ctx, `
        INSERT INTO transactions (user_id, amount, operation_type, description, related_user_id)
        VALUES ($1, $2, $3, $4, $5)
    `, transaction.UserID, transaction.Amount, transaction.OperationType, transaction.Description, transaction.RelatedUserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *transactionRepositoryImpl) GetLastTransactions(ctx context.Context, userID int64, limit int) ([]entity.Transaction, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, user_id, amount, operation_type, description, created_at, related_user_id
        FROM transactions
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []entity.Transaction
	for rows.Next() {
		var transaction entity.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Amount,
			&transaction.OperationType,
			&transaction.Description,
			&transaction.CreatedAt,
			&transaction.RelatedUserID,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
