package usecase

import (
	"context"
	"fmt"

	"finances/internal/entity"
	"finances/internal/repository"

	"github.com/jackc/pgx/v4"
)

type UserService interface {
	Deposit(ctx context.Context, userID int64, amount float64) error
	Transfer(ctx context.Context, fromUserID, toUserID int64, amount float64) error
	GetLastTransactions(ctx context.Context, userID int64, limit int) ([]entity.Transaction, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepository
	txRepo   repository.TransactionRepository
	db       *pgx.Conn
}

func NewUserService(userRepo repository.UserRepository, txRepo repository.TransactionRepository, db *pgx.Conn) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
		txRepo:   txRepo,
		db:       db,
	}
}

func (s *userServiceImpl) Deposit(ctx context.Context, userID int64, amount float64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = s.userRepo.UpdateBalance(ctx, tx, userID, amount)
	if err != nil {
		return err
	}

	transaction := entity.Transaction{
		UserID:        userID,
		Amount:        amount,
		OperationType: "deposit",
		Description:   fmt.Sprintf("Deposited %.2f", amount),
	}
	err = s.txRepo.CreateTransaction(ctx, tx, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) Transfer(ctx context.Context, fromUserID, toUserID int64, amount float64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	fromBalance, err := s.userRepo.GetBalance(ctx, tx, fromUserID)
	if err != nil {
		return err
	}
	if fromBalance < amount {
		return fmt.Errorf("insufficient funds")
	}

	err = s.userRepo.UpdateBalance(ctx, tx, fromUserID, -amount)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdateBalance(ctx, tx, toUserID, amount)
	if err != nil {
		return err
	}

	transaction := entity.Transaction{
		UserID:        fromUserID,
		Amount:        -amount,
		OperationType: "transfer",
		Description:   fmt.Sprintf("Transferred %.2f to user %d", amount, toUserID),
		RelatedUserID: toUserID,
	}
	err = s.txRepo.CreateTransaction(ctx, tx, transaction)
	if err != nil {
		return err
	}

	transaction = entity.Transaction{
		UserID:        toUserID,
		Amount:        amount,
		OperationType: "transfer",
		Description:   fmt.Sprintf("Received %.2f from user %d", amount, fromUserID),
		RelatedUserID: fromUserID,
	}
	err = s.txRepo.CreateTransaction(ctx, tx, transaction)
	if err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) GetLastTransactions(ctx context.Context, userID int64, limit int) ([]entity.Transaction, error) {
	return s.txRepo.GetLastTransactions(ctx, userID, limit)
}
