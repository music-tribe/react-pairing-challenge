package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	client *mongo.Client
	logger MongoDatabaseLogger
}

type MongoDatabaseLogger interface {
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

func OpenMongoConnection(url string, logger MongoDatabaseLogger) (*MongoDatabase, error) {
	cli, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		logger.Errorf("OpenMongoConnection: failed to open new mgo.Session >> %v", err)
		return nil, err
	}

	if err := cli.Ping(context.Background(), nil); err != nil {
		logger.Errorf("OpenMongoConnection: failed on ping to connection >> %v", err)
		return nil, err
	}

	logger.Infof("OpenMongoConnection: MongoDB Connected")

	return &MongoDatabase{
		client: cli,
		logger: logger,
	}, nil
}

func (mdb *MongoDatabase) CloseMongoConnection() error {
	mdb.logger.Infof("CloseMongoConnection: MongoDB disconnected")
	return mdb.client.Disconnect(context.TODO())
}
