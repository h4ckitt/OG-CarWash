package repository

import "car_wash/model"

type Repo interface {
	RegisterCarWash(userID string, carWash model.CarWash) (string, error)
	RegisterOwner(owner model.Owner) error
	FetchCarWashDataByDate(date string, carWashID string) (model.WebSocketResult, error)
	FetchAllCarWashDataByDate(date string, ownerID string) (model.WebSocketResult, error)
	VerifyAPIKey(key string) (string, error)
	RetrieveAPIKey(uid string) (string, error)
	StoreAPIKey(uid string, key string) error
	UpdateAPIKey(uid string, key string) error
	SaveWashDetails(uid string, wash model.Wash) error
}
