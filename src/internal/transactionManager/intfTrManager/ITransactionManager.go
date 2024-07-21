package intfTrManager

import "context"

//go:generate mockgen -source=ITransactionManager.go -destination=../../tests/unitTests/mocks/mockITransactionManager.go --package=mocks

type ITransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
