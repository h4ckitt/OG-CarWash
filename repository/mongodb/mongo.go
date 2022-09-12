package mongodb

import (
	"car_wash/config"
	"car_wash/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

func (r Repo) RegisterOwner(owner model.Owner) error {
	return nil
}

func (r Repo) FetchDataByDate(date string) (model.WebSocketResult, error) {
	//TODO implement me
	clientResults := make(map[string][]model.CarWashes)
	carWashResults := make(map[int][]model.Wash)

	savedCars := r.database.Collection("saved_cars")

	res, err := savedCars.Find(context.Background(), bson.D{{"Day entered", bson.D{{"$eq", date}}}})

	if err != nil {
		return model.WebSocketResult{}, err
	}

	defer func() { _ = res.Close(context.Background()) }()

	for res.Next(context.Background()) {
		var w model.Wash

		if err = res.Decode(&w); err != nil {
			return model.WebSocketResult{}, err
		}

		id := w.CarWashID

		// Remove Fields We Don't Need From The Output JSON
		w.DateEntered = ""
		w.CarWashID = 0

		carWashResults[id] = append(carWashResults[id], w)
	}

	for id, washes := range carWashResults {
		clientResults[washes[0].ClientNumber] = append(clientResults[washes[0].ClientNumber], model.CarWashes{
			CarWashID:   id,
			CarsEntered: len(washes),
			Cars:        washes,
		})
	}

	var result model.WebSocketResult

	result.Date = date

	for clients, wash := range clientResults {
		result.Clients = append(result.Clients, model.WebSocketClientResult{
			ClientNumber: clients,
			Result:       wash,
		})
	}
	return result, nil
}
