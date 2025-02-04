package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type userRepositoryImpl struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) GetBalance(ctx context.Context, tx pgx.Tx, userID int64) (float64, error) {
	var balance float64
	err := tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *userRepositoryImpl) UpdateBalance(ctx context.Context, tx pgx.Tx, userID int64, amount float64) error {
	_, err := tx.Exec(ctx, "UPDATE users SET balance = balance + $1 WHERE id = $2", amount, userID)
	if err != nil {
		return err
	}
	return nil
}
