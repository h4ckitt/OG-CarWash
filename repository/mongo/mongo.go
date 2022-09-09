package mongo

import (
	"car_wash/config"
	"car_wash/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoClient() (*Repo, error) {
	var uri string
	mongoConfig := config.GetConfig().MongoConfig

	uri = fmt.Sprintf("mongodb://%s:%s", mongoConfig.DatabaseHost, mongoConfig.DatabasePort)

	ctx := context.Background()

	opt := options.Client().ApplyURI(uri)

	if mongoConfig.DatabaseUser != "" {
		opt.SetAuth(options.Credential{Username: mongoConfig.DatabaseUser, Password: mongoConfig.DatabasePassword})
	}

	client, err := mongo.Connect(ctx, opt)
	db := client.Database(mongoConfig.DatabaseName)

	if err != nil {
		return nil, err
	}

	return &Repo{
		client:   client,
		database: db,
	}, nil
}

func (r Repo) RegisterOwner(owner model.Owner) {
	//TODO implement me
	panic("implement me")
}

func (r Repo) FetchDataByDate(date string) model.WebSocketResult {
	//TODO implement me
	panic("implement me")
}
