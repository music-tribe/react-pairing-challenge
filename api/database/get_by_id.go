package database

import (
	"context"

	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mdb *MongoDatabase) GetById(featureId uuid.UUID) (*domain.Feature, error) {
	coll := mdb.client.Database("pair-challenge").Collection("features")

	q := coll.FindOne(context.Background(), bson.M{"_id": featureId})

	t := new(domain.Feature)
	err := q.Decode(t)
	if err != nil {
		mdb.logger.Errorf("database.GetById: mongo.Decode >> %v", err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return t, nil
}
