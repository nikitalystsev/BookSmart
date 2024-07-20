package interfaces

import "context"

type ITransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
