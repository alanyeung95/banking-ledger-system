package accounts

import "context"

// Repository is the account repo
type Repository interface {
	Upsert(ctx context.Context, id string, account Account) (*Account, error)
}
