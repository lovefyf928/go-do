package mongo

import (
	"context"
	"fmt"
	"go-do/common/conf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongodb *mongo.Database

func init() {
	clientOption := options.Client().ApplyURI(conf.ConfigInfo.DataSource.Mongo.Uri)
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	Mongodb = client.Database(conf.ConfigInfo.DataSource.Mongo.Db)

	fmt.Println("Connected to MongoDB!")
}
