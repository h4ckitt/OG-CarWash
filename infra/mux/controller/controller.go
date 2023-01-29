package controller

import (
	"car_wash/apperror"
	"car_wash/infra/mux/helper"
	"car_wash/model"
	"car_wash/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Controller struct {
	Service service.Svc
}

var upgrader = websocket.Upgrader{}

func NewController(svc service.Svc) *Controller {
	return &Controller{svc}
}

func (controller *Controller) CheckHash(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]

	if hash == "" {
		helper.ReturnFailure(w, &apperror.BadRequest)
		return
	}

	key, err := controller.Service.CheckCreds(hash)

	if err != nil {
		helper.ReturnFailure(w, err)
		return
	}

	helper.ReturnSuccess(w, map[string]string{"apiKey": key})
}

func (controller *Controller) AddCredsToJar(w http.ResponseWriter, r *http.Request) {

	h := struct {
		Hash string `json:"hash"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		helper.ReturnFailure(w, &apperror.BadRequest)
		return
	}

	err := controller.Service.CacheCreds(r.Context(), h.Hash)

	if err != nil {
		helper.ReturnFailure(w, err)
	}
}

func (controller *Controller) RegisterNewOwner(w http.ResponseWriter, r *http.Request) {
	var owner model.Owner

	if err := json.NewDecoder(r.Body).Decode(&owner); err != nil {
		helper.ReturnFailure(w, &apperror.BadRequest)
		return
	}

	key, err := controller.Service.RegisterNewOwner(r.Context(), owner)

	if err != nil {
		helper.ReturnFailure(w, err)
		return
	}

	helper.ReturnSuccess(w, map[string]string{"apiKey": key})

	//w.WriteHeader(http.StatusCreated)
}

func (controller *Controller) UpgradeWss(w http.ResponseWriter, r *http.Request) {
	updatesChan := controller.Service.GetUpdatesChannel(r.Context())

	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Upgrade: ", err)
		return
	}

	defer func() { _ = c.Close() }()

	go func() {
		for {
			_ = c.WriteJSON(<-updatesChan)
		}
	}()

	for {
		_, message, err := c.ReadMessage()

		if err != nil {
			log.Println("Read: ", err)
			break
		}

		//log.Printf("Received: %s", message)

		if match, err := regexp.MatchString("^(?:[0-9]{2}.){2}[0-9]{4}$", string(message)); err != nil || !match {
			continue
		}

		res, err := controller.Service.FetchDataByDate(r.Context(), string(message))

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

func (controller *Controller) RegisterCarWash(w http.ResponseWriter, r *http.Request) {
	var carwash model.CarWash

	if err := json.NewDecoder(r.Body).Decode(&carwash); err != nil {
		helper.ReturnFailure(w, &apperror.BadRequest)
		return
	}

	id, err := controller.Service.RegisterCarWash(r.Context(), carwash)

	if err != nil {
		helper.ReturnFailure(w, err)
		return
	}

	helper.ReturnSuccess(w, map[string]string{"message": "created successfully", "id": id})
}

func (controller *Controller) RegisterWash(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		helper.ReturnFailure(w, &apperror.BadRequest)
		return
	}

	wash := model.Wash{
		CarWashID:   r.FormValue("carWashID"),
		NumberPlate: r.FormValue("license"),
		DateEntered: r.FormValue("dateEntered"),
	}

	file, header, err := r.FormFile("image")

	if err != nil {
		helper.ReturnFailure(w, &apperror.BadRequest)
		return
	}

	fileBits := strings.Split(header.Filename, ".")

	wash.ImageExt = fileBits[len(fileBits)-1]
	wash.Image = file

	err = controller.Service.SaveWashDetails(r.Context(), wash)

	if err != nil {
		helper.ReturnFailure(w, err)
		return
	}
}

//func (controller *Controller)
