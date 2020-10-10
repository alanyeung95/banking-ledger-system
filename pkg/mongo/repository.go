package mongo

import (
	"context"
	"fmt"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)



// NewClient is the function to create an new mongo connection
func NewClient(addresses, username, password, database string) (*mongo.Client, error) {
	// 	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s", username, password, addresses, database)
	uri := fmt.Sprintf("mongodb://%s/%s", addresses, database)
	opts := options.Client().
		ApplyURI(uri)
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

func newCollection(client *mongo.Client, database, collection string) (*mongo.Collection, error) {
	c := client.Database(database).Collection(collection)
	return c, nil
}

// upsert create or update the target document and return created or updated document
func upsert(ctx context.Context, c *mongo.Collection, filter interface{}, update interface{}, result interface{}) error {
	err := c.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetUpsert(true),
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func findOne(ctx context.Context, c *mongo.Collection, filter interface{}, result interface{}) error {
	err := c.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		return err
	}

	return nil
}

// findAndUpdateOne update the target document and return updated document
func findAndUpdateOne(ctx context.Context, c *mongo.Collection, filter interface{}, update interface{}, result interface{}) error {
	err := c.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		return err
	}
	return nil
}

func findAll(ctx context.Context, c *mongo.Collection, filter interface{}, results interface{}, opts ...*options.FindOptions) error {
	cursor, err := c.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}

	return cursor.All(ctx, results)
}