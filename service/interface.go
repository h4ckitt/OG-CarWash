package service

import (
	"car_wash/model"
	"context"
)

type Svc interface {
	SaveWashDetails(ctx context.Context, wash model.Wash) error
	RegisterCarWash(ctx context.Context, carWash model.CarWash) (string, error)
	RegisterNewOwner(ctx context.Context, owner model.Owner) (string, error)
	FetchDataByDate(ctx context.Context, date string) (model.WebSocketResult, error)
	GetUpdatesChannel(ctx context.Context) <-chan model.Wash
	CacheCreds(ctx context.Context, hash string) error
	CheckCreds(hash string) (string, error)
}
