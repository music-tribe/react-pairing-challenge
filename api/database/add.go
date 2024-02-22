package database

import (
	"context"
	"errors"

	"github.com/music-tribe/react-pairing-challenge/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mdb *MongoDatabase) Add(task *domain.Task) error {
	coll := mdb.client.Database("pair-challenge").Collection("tasks")

	b, err := bson.Marshal(task)
	if err != nil {
		return err
	}

	_, err = coll.InsertOne(context.Background(), b)
	if err != nil {
		wrEx := mongo.WriteException{}
		if errors.As(err, &wrEx) {
			if wrEx.HasErrorCode(11000) {
				return ErrDuplicate
			}
		}

		return err
	}

	return nil
}
