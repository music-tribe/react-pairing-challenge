package database

import (
	"context"
	"errors"

	"github.com/music-tribe/react-pairing-challenge/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mdb *MongoDatabase) Add(feature *domain.Feature) error {
	coll := mdb.client.Database("pair-challenge").Collection("features")

	b, err := bson.Marshal(feature)
	if err != nil {
		mdb.logger.Errorf("database.Add: bson.Marshal >> %v", err)
		return err
	}

	_, err = coll.InsertOne(context.Background(), b)
	if err != nil {
		mdb.logger.Errorf("database.Add: mongo.InsertOne >> %v", err)
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
