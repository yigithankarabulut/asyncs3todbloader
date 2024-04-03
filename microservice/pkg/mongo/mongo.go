package mongo

import (
	"context"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var DB *mongo.Database

func Connect(dbConfig config.Database) error {
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
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	db := client.Database(dbConfig.Name)
	DB = db
	return nil
}
