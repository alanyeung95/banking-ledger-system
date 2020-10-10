package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alanyeung95/banking-ledger-system/pkg/accounts"

)

// NewAccountRepository is the repo to store account model
func NewAccountRepository(client *mongo.Client, database, collection string) (*AccountRepository, error) {
	c, err := newCollection(client, database, collection)
	if err != nil {
		return nil, err
	}

	return &AccountRepository{c}, nil
}

type AccountRepository struct {
	collection *mongo.Collection
}

// interface check
var _ accounts.Repository = (*AccountRepository)(nil)

// Upsert returns the Account record being successfully created or updated
func (r *AccountRepository) Upsert(ctx context.Context, id string, Account accounts.Account) (*accounts.Account, error) {
	var result accounts.Account
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": Account,
	}
	if err := upsert(ctx, r.collection, filter, update, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Find returns the Account record being successfully created or updated
func (r *AccountRepository) Find(ctx context.Context, id string) (*accounts.Account, error) {
	var result accounts.Account
	filter := bson.M{"_id": id}
	if err := r.collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &result, nil
}

// Update returns the Account record being successfully updated
func (r *AccountRepository) UpdateBalance(ctx context.Context, id string, amount int) (*accounts.Account, error) {
	var result accounts.Account
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{
			"balance": amount,
		},
	}
	
	err :=  findAndUpdateOne(ctx, r.collection, filter, update, &result)
	if err != nil {
		return nil, err
	}

	return &result,nil
}