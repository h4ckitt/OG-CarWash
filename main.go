package main

import (
	"car_wash/config"
	"car_wash/infra/auth/authenticator"
	"car_wash/infra/mux/controller"
	"car_wash/infra/mux/router"
	"car_wash/repository/postgres"
	"car_wash/service"
	"fmt"
	"log"
	"net/http"
)

func main() {

	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}

	//repo, err := mongodb.NewMongoClient()
	repo, err := postgres.NewPostgresRepo()

	if err != nil {
		log.Fatalln(err)
	}

	auth, err := authenticator.NewAuthenticator(repo)

	if err != nil {
		log.Fatalln(err)
	}

	ctrl := controller.NewController(service.NewService(repo))

	r := router.InitRouter(ctrl, auth)

	log.Println("Starting Server ......")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.GetConfig().RunConfig.Port), r); err != nil {
		log.Panicln(err)
	}

}
