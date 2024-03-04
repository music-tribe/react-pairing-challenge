package database

import (
	"context"

	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mdb *MongoDatabase) Get(userId, id uuid.UUID) (*domain.Feature, error) {
	coll := mdb.client.Database("pair-challenge").Collection("features")

	q := coll.FindOne(context.Background(), bson.M{"_id": id, "userId": userId})

	t := new(domain.Feature)
	err := q.Decode(t)
	if err != nil {
		mdb.logger.Errorf("database.Get: mongo.Decode >> %v", err)
		if err == mongo.ErrNoDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return t, nil
}
