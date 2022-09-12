package service

import (
	"car_wash/model"
	"car_wash/repository"
	"errors"
	"log"
)

type CarWashSvc struct {
	Repo repository.Repo
}

func NewService(repository repository.Repo) *CarWashSvc {
	return &CarWashSvc{repository}
}

func (c CarWashSvc) RegisterNewOwner(owner model.Owner) error {
	err := c.Repo.RegisterOwner(owner)

	if err != nil {
		log.Println(err)
	}

	return errors.New("an error occurred while trying to register")
}

func (c CarWashSvc) FetchDataByDate(date string) model.WebSocketResult {
	res, err := c.Repo.FetchDataByDate(date)

	if err != nil {
		log.Println(err)
	}

	return res
}
