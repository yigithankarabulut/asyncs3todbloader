package mongo

import (
	"context"
	"github.com/yigithankarabulut/asyncs3todbloader/job/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func ConnectMongo(dbConfig config.Database) (*mongo.Database, error) {
	var connString string
	if dbConfig.User == "" && dbConfig.Pass == "" {
		connString = "mongodb://" + dbConfig.Host + ":" + dbConfig.Port
	} else {
		connString = "mongodb://" + dbConfig.User + ":" + dbConfig.Pass + "@" + dbConfig.Host + ":" + dbConfig.Port
	}
	clientOptions := options.Client().ApplyURI(connString)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	db := client.Database(dbConfig.Name)
	log.Println("Connected to MongoDB")
	_ = CreateCollection(db, dbConfig.ProductCollection)
	_ = CreateCollection(db, dbConfig.ObjectInfoCollection)
	return db, nil
}

func CreateCollection(db *mongo.Database, collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}
