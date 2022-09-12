package controller

import (
	"car_wash/model"
	"car_wash/service"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"regexp"
)

type Controller struct {
	Service service.Svc
}

var upgrader = websocket.Upgrader{}

func NewController(svc service.Svc) *Controller {
	return &Controller{svc}
}

func (controller *Controller) RegisterNewOwner(w http.ResponseWriter, r *http.Request) {
	var owner model.Owner

	if err := json.NewDecoder(r.Body).Decode(&owner); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "bad request received"})
	}

	err := controller.Service.RegisterNewOwner(owner)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (controller *Controller) UpgradeWss(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Upgrade: ", err)
		return
	}

	defer func() { _ = c.Close() }()

	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Println("Read: ", err)
			break
		}

		log.Printf("Received: %s", message)

		if match, err := regexp.MatchString("^(?:[0-9]{2}.){2}[0-9]{4}$", string(message)); err != nil || !match {
			continue
		}

		res, err := controller.Service.FetchDataByDate(string(message))

		if err != nil {
			_ = c.WriteJSON(map[string]string{"message": "an error occurred, please try again"})
			continue
		}

		if err = c.WriteJSON(res); err != nil {
			log.Println(err)
			continue
		}
	}
}
