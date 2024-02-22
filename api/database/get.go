package database

import (
	"context"

	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (mdb *MongoDatabase) Get(id uuid.UUID) (*domain.Task, error) {
	coll := mdb.client.Database("pair-challenge").Collection("tasks")

	q := coll.FindOne(context.Background(), bson.M{"_id": id})

	t := new(domain.Task)
	err := q.Decode(t)
	if err != nil {
		mdb.logger.Errorf("database.Get: mongo.Decode >> %v", err)
		return nil, err
	}

	return t, nil
}
