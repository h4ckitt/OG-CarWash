package service

import (
	"car_wash/model"
	"car_wash/repository"
)

type CarWashSvc struct {
	Repo repository.Repo
}

func NewService(repository repository.Repo) *CarWashSvc {
	return &CarWashSvc{repository}
}

func (c CarWashSvc) RegisterNewOwner(owner model.Owner) {
	//TODO implement me
	panic("implement me")
}

func (c CarWashSvc) FetchDataByDate(date string) model.WebSocketResult {
	//TODO implement me
	panic("implement me")
}
