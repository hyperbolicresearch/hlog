package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func Client(uri string) *mongo.Client {
	if uri == "" {
		panic(errors.New("invalid mongodb uri"))
	}
	clientOptions := options.Client().ApplyURI(uri)
	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	return client
	
}
