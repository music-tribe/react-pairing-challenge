package database

import (
	"context"

	"github.com/music-tribe/react-pairing-challenge/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mdb *MongoDatabase) Update(task *domain.Task) error {
	coll := mdb.client.Database("pair-challenge").Collection("tasks")

	q := coll.FindOneAndReplace(context.Background(), bson.M{"_id": task.Id, "userId": task.UserId}, task)

	if err := q.Err(); err != nil {
		mdb.logger.Errorf("database.Update: mongo.Err >> %v", err)
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return err
	}

	return nil
}
