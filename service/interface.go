package service

import "car_wash/model"

type Svc interface {
	RegisterNewOwner(owner model.Owner) error
	FetchDataByDate(date string) (model.WebSocketResult, error)
}
