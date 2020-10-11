package transactions

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Service interface
type Service interface {
	RecordTransaction(ctx context.Context, r *http.Request, transcation *Transaction) (*Transaction, error)
	GetTransactionsByID(ctx context.Context, r *http.Request, id string, asc int) ([]Transaction, error)
	GetTransactionByID(ctx context.Context, r *http.Request, id string) (*Transaction, error)
	DeleteTransactions(ctx context.Context, r *http.Request, id string) error
}

type service struct {
	repository Repository
}

// NewService start the new service
func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) RecordTransaction(ctx context.Context, r *http.Request, transaction *Transaction) (*Transaction, error) {
	transaction.ID = uuid.NewV4().String()
	return s.repository.Upsert(ctx, transaction.ID, *transaction)
}

func (s *service) GetTransactionsByID(ctx context.Context, r *http.Request, id string, asc int) ([]Transaction, error) {
	return s.repository.FindAll(ctx, id, asc)
}

func (s *service) GetTransactionByID(ctx context.Context, r *http.Request, id string) (*Transaction, error) {
	return s.repository.Find(ctx, id)
}

func (s *service) DeleteTransactions(ctx context.Context, r *http.Request, id string)  error {
	return s.repository.Delete(ctx, id)
}