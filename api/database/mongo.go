package database

import (
	"context"
	"errors"

	"github.com/music-tribe/react-pairing-challenge/domain"
	"github.com/music-tribe/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

func (mdb *MongoDatabase) Get(id uuid.UUID) (*domain.Task, error) {
	coll := mdb.client.Database("pair-challenge").Collection("tasks")

	q := coll.FindOne(context.Background(), bson.M{"_id": id})

	t := new(domain.Task)
	err := q.Decode(t)

	return t, err
}
