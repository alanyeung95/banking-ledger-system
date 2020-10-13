package transactions

import "context"

// Repository is the transaction repo
type Repository interface {
	Upsert(ctx context.Context, id string, transaction Transaction) (*Transaction, error)
	FindAll(ctx context.Context, accountID string, asc int) ([]Transaction, error)
	Find(ctx context.Context, id string) (*Transaction, error)
	Delete(ctx context.Context, id string) error
}
