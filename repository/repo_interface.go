package repository

import "car_wash/model"

type Repo interface {
	RegisterOwner(owner model.Owner) error
	FetchDataByDate(date string) (model.WebSocketResult, error)
}
