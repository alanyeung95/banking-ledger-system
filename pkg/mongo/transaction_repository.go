package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

// Upsert returns the Transaction records that matching the criteria
func (r *TransactionRepository) FindAll(ctx context.Context, id string, asc int) ([]transactions.Transaction, error) {
	var transactionList []transactions.Transaction
	filter := bson.M{
		"$or": bson.A{
			bson.M{"body.from": id},
			bson.M{"body.to": id},
		},
	}

	findOption := options.Find()
	findOption.SetSort(bson.M{"time": asc})

	if err := findAll(ctx, r.collection, filter, &transactionList, findOption); err != nil{
		return nil, err
	}
	return transactionList, nil
}

// Find returns the Transaction record being successfully created or updated
func (r *TransactionRepository) Find(ctx context.Context, id string) (*transactions.Transaction, error) {
	var result transactions.Transaction
	filter := bson.M{"_id": id}
	if err := r.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &result, nil
}
