package interfaces

import "context"

type ITransactionManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
