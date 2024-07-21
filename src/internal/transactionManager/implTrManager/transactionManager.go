package implTrManager

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type TransactionManager struct {
	transactionManager manager.Manager
}

func NewTransactionManager(transactionManager manager.Manager) *TransactionManager {
	return &TransactionManager{transactionManager: transactionManager}
}

func (trm *TransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return trm.transactionManager.Do(ctx, fn)
}
