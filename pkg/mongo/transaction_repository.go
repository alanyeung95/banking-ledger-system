package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alanyeung95/banking-ledger-system/pkg/transactions"

)

// NewTransactionRepository is the repo to store transaction model
func NewTransactionRepository(client *mongo.Client, database, collection string) (*TransactionRepository, error) {
	c, err := newCollection(client, database, collection)
	if err != nil {
		return nil, err
	}

	return &TransactionRepository{c}, nil
}

type TransactionRepository struct {
	collection *mongo.Collection
}

// interface check
var _ transactions.Repository = (*TransactionRepository)(nil)

// Upsert returns the Transaction record being successfully created or updated
func (r *TransactionRepository) Upsert(ctx context.Context, id string, Transaction transactions.Transaction) (*transactions.Transaction, error) {
	var result transactions.Transaction
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": Transaction,
	}
	if err := upsert(ctx, r.collection, filter, update, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
