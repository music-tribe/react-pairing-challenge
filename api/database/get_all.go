package database

import (
	"context"

	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mdb *MongoDatabase) GetAll(userId uuid.UUID) ([]*domain.Feature, error) {
	coll := mdb.client.Database("pair-challenge").Collection("features")

	q, err := coll.Find(context.Background(), bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}

	ts := make([]*domain.Feature, 0)

	if err = q.All(context.Background(), &ts); err != nil {
		mdb.logger.Errorf("database.GetAll: mongo.All >> %v", err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return ts, nil
}
