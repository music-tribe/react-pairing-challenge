package database

import (
	"context"

	"github.com/music-tribe/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (mdb *MongoDatabase) Delete(taskId uuid.UUID) error {
	coll := mdb.client.Database("pair-challenge").Collection("tasks")

	res, err := coll.DeleteOne(context.Background(), bson.M{"_id": taskId})
	if err != nil {
		mdb.logger.Errorf("database.Delete: mongo.DeleteOne >> %v", err)
		return err
	}

	if res.DeletedCount == 0 {
		mdb.logger.Errorf("database.Delete: mongo.DeleteOne >> %v", ErrNotFound)
		return ErrNotFound
	}

	return nil
}
