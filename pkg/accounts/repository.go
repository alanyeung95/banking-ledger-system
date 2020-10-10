package accounts

import "context"

// Repository is the account repo
type Repository interface {
	Upsert(ctx context.Context, id string, account Account) (*Account, error)
	Find(ctx context.Context, id string) (*Account, error)
	UpdateBalance(ctx context.Context, id string, amount int ) (*Account, error)
}
