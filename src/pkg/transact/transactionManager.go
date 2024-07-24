package transact

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

//go:generate mockgen -source=ITransactionManager.go -destination=../../tests/unitTests/mocks/mockITransactionManager.go --package=mocks

type ITransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionManager struct {
	transactionManager manager.Manager
}

func NewTransactionManager(transactionManager manager.Manager) *TransactionManager {
	return &TransactionManager{transactionManager: transactionManager}
}

func (trm *TransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return trm.transactionManager.Do(ctx, fn)
}
