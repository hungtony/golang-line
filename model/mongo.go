package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type MongoRepository struct {
	client *mongo.Client
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{
		client: client,
	}
}

func (r *MongoRepository) SaveMessage(ctx context.Context, message *Message) error {

	log.Println(message)
	collection := r.client.Database("local").Collection("user_info")

	_, err := collection.InsertOne(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) GetMessagesByUserID(ctx context.Context, userID string) ([]*Message, error) {
	collection := r.client.Database("line_bot").Collection("messages")

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var messages []*Message
	for cursor.Next(ctx) {
		message := &Message{}
		err := cursor.Decode(message)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}
