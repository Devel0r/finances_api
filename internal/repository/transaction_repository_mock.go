package repository

import (
	"context"
	"finances/internal/entity"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
)

type TransactionRepositoryMock struct {
	mock.Mock
}

func (m *TransactionRepositoryMock) CreateTransaction(ctx context.Context, tx pgx.Tx, transaction entity.Transaction) error {
	args := m.Called(ctx, tx, transaction)
	return args.Error(0)
}

func (m *TransactionRepositoryMock) GetLastTransactions(ctx context.Context, userID int64, limit int) ([]entity.Transaction, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]entity.Transaction), args.Error(1)
}
