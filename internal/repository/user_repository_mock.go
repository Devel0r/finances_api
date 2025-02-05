package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetBalance(ctx context.Context, tx pgx.Tx, userID int64) (float64, error) {
	args := m.Called(ctx, tx, userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *UserRepositoryMock) UpdateBalance(ctx context.Context, tx pgx.Tx, userID int64, amount float64) error {
	args := m.Called(ctx, tx, userID, amount)
	return args.Error(0)
}
