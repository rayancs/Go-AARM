package repo

import (
	"app/logger"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	db *mongo.Database
}

func NewMongoDB(uri, dbName string) *Mongo {
	lg := logger.New()
	// create a context and a cancelation token so there is a timout
	ctx, cancelToken := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelToken() // dispose cancelToken at end of func
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri))
	if err != nil {
		lg.Error(err.Error(), err)
	}
	db := client.Database(dbName)
	return &Mongo{
		db: db,
	}
}

func (m *Mongo) fromCollection(collName string) *mongo.Collection {
	return m.db.Collection(collName)
}
