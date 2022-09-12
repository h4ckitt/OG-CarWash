package service

import (
	"car_wash/model"
	"car_wash/repository"
	"errors"
	"log"
	"time"
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

func (c CarWashSvc) FetchDataByDate(date string) (model.WebSocketResult, error) {

	t, err := time.Parse("02.01.2006", date)

	if err != nil {
		return model.WebSocketResult{}, errors.New("invalid date specified")
	}

	res, err := c.Repo.FetchDataByDate(t.Format("2006-01-02"))

	if err != nil {
		log.Println(err)
	}

	return res, err
}
