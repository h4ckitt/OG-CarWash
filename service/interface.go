package service

import "car_wash/model"

type Svc interface {
	RegisterNewOwner(owner model.Owner)
	FetchDataByDate(date string) model.WebSocketResult
}
