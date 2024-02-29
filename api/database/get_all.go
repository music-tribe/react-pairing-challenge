package database

import (
	"context"

	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mdb *MongoDatabase) GetAll(userId uuid.UUID) ([]*domain.Task, error) {
	coll := mdb.client.Database("pair-challenge").Collection("tasks")

	q, err := coll.Find(context.Background(), bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}

	ts := make([]*domain.Task, 0)

	if err = q.All(context.Background(), &ts); err != nil {
		mdb.logger.Errorf("database.Get: mongo.Decode >> %v", err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return ts, nil
}
