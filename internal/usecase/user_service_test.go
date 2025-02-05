package usecase

import (
	"context"
	"testing"

	"finances/internal/entity"
	"finances/internal/repository"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

func TestDeposit(t *testing.T) {
	userRepo := new(repository.UserRepositoryMock)
	txRepo := new(repository.TransactionRepositoryMock)
	db, _ := pgx.Connect(context.Background(), "postgresql://user:pass@localhost/db") // Можно использовать фиктивное подключение
	userService := NewUserService(userRepo, txRepo, db)

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when starting a transaction", err)
	}
	defer tx.Rollback(ctx)

	userRepo.On("UpdateBalance", ctx, tx, int64(1), float64(100.0)).Return(nil)
	txRepo.On("CreateTransaction", ctx, tx, entity.Transaction{
		UserID:        1,
		Amount:        100.0,
		OperationType: "deposit",
		Description:   "Deposited 100.00",
	}).Return(nil)

	err = userService.Deposit(ctx, 1, 100.0)
	assert.NoError(t, err)

	userRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestTransfer(t *testing.T) {
	userRepo := new(repository.UserRepositoryMock)
	txRepo := new(repository.TransactionRepositoryMock)
	db, _ := pgx.Connect(context.Background(), "postgresql://user:pass@localhost/db") // Можно использовать фиктивное подключение
	userService := NewUserService(userRepo, txRepo, db)

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when starting a transaction", err)
	}
	defer tx.Rollback(ctx)

	userRepo.On("GetBalance", ctx, tx, int64(1)).Return(float64(200.0), nil)
	userRepo.On("UpdateBalance", ctx, tx, int64(1), float64(-100.0)).Return(nil)
	userRepo.On("UpdateBalance", ctx, tx, int64(2), float64(100.0)).Return(nil)
	txRepo.On("CreateTransaction", ctx, tx, entity.Transaction{
		UserID:        1,
		Amount:        -100.0,
		OperationType: "transfer",
		Description:   "Transferred 100.00 to user 2",
		RelatedUserID: 2,
	}).Return(nil)
	txRepo.On("CreateTransaction", ctx, tx, entity.Transaction{
		UserID:        2,
		Amount:        100.0,
		OperationType: "transfer",
		Description:   "Received 100.00 from user 1",
		RelatedUserID: 1,
	}).Return(nil)

	err = userService.Transfer(ctx, 1, 2, 100.0)
	assert.NoError(t, err)

	userRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestTransfer_InsufficientFunds(t *testing.T) {
	userRepo := new(repository.UserRepositoryMock)
	txRepo := new(repository.TransactionRepositoryMock)
	db, _ := pgx.Connect(context.Background(), "postgresql://user:pass@localhost/db") // Можно использовать фиктивное подключение
	userService := NewUserService(userRepo, txRepo, db)

	ctx := context.Background()
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when starting a transaction", err)
	}
	defer tx.Rollback(ctx)

	userRepo.On("GetBalance", ctx, tx, int64(1)).Return(float64(50.0), nil)

	err = userService.Transfer(ctx, 1, 2, 100.0)
	assert.EqualError(t, err, "insufficient funds")

	userRepo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestGetLastTransactions(t *testing.T) {
	userRepo := new(repository.UserRepositoryMock)
	txRepo := new(repository.TransactionRepositoryMock)
	db, _ := pgx.Connect(context.Background(), "postgresql://user:pass@localhost/db") // Можно использовать фиктивное подключение
	userService := NewUserService(userRepo, txRepo, db)

	ctx := context.Background()

	transactions := []entity.Transaction{
		{
			ID:            1,
			UserID:        1,
			Amount:        100.0,
			OperationType: "deposit",
			Description:   "Deposited 100.00",
			CreatedAt:     "2023-01-01T00:00:00Z",
			RelatedUserID: 0,
		},
		{
			ID:            2,
			UserID:        1,
			Amount:        -50.0,
			OperationType: "transfer",
			Description:   "Transferred 50.00 to user 2",
			CreatedAt:     "2023-01-02T00:00:00Z",
			RelatedUserID: 2,
		},
	}

	txRepo.On("GetLastTransactions", ctx, int64(1), 10).Return(transactions, nil)

	result, err := userService.GetLastTransactions(ctx, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	txRepo.AssertExpectations(t)
}
