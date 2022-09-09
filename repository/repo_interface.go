package repository

import "car_wash/model"

type Repo interface {
	RegisterOwner(owner model.Owner)
	FetchDataByDate(date string) model.WebSocketResult
}
