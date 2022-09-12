package main

import (
	"car_wash/config"
	"car_wash/infra/mux/controller"
	"car_wash/infra/mux/router"
	"car_wash/repository/mongodb"
	"car_wash/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
		return
	}

	if match, err := regexp.MatchString("^[0-9]+$", os.Args[1]); err != nil || !match {
		log.Fatalln("invalid port specified")
	}

	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}

	repo, err := mongodb.NewMongoClient()

	if err != nil {
		log.Fatalln(err)
	}

	ctrl := controller.NewController(service.NewService(repo))

	log.Println("Starting Server ......")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Args[1]), router.InitRouter(ctrl)); err != nil {
		log.Panicln(err)
	}

}
