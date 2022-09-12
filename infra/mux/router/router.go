package router

import (
	"car_wash/infra/mux/controller"
	"github.com/gorilla/mux"
	"net/http"
)

func InitRouter(controller *controller.Controller) *mux.Router {
	router := mux.NewRouter()
	pathSubrouter := router.PathPrefix("/carwash").Subrouter()

	pathSubrouter.HandleFunc("", controller.UpgradeWss).Methods(http.MethodGet)
	pathSubrouter.HandleFunc("/registerOwner", controller.RegisterNewOwner).Methods(http.MethodPost)

	return router
}
