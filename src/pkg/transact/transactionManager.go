package transact

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

//go:generate mockgen -source=transactionManager.go -destination=../../internal/tests/unitTests/mocks/mockTransactionManager.go --package=mocks

type ITransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type TransactionManager struct {
	transactionManager manager.Manager
}

func NewTransactionManager(transactionManager manager.Manager) *TransactionManager {
	return &TransactionManager{transactionManager: transactionManager}
}

func (trm *TransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return trm.transactionManager.Do(ctx, fn)
}
