package postgres

import (
	"car_wash/config"
	"car_wash/model"
	"database/sql"
	"fmt"
)

type Repo struct {
	conn *sql.DB
}

func NewPostgresRepo() (*Repo, error) {
	pgConfig := config.GetConfig().PGConfig

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", pgConfig.DatabaseHost, pgConfig.DatabasePort, pgConfig.DatabaseUser, pgConfig.DatabasePassword, pgConfig.DatabaseName))

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Repo{
		conn: db,
	}, nil
}

func (r Repo) RegisterOwner(owner model.Owner) {
	//TODO implement me
	panic("implement me")
}

func (r Repo) FetchDataByDate(date string) model.WebSocketResult {
	//TODO implement me
	panic("implement me")
}
