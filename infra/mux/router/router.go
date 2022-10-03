package router

import (
	"car_wash/infra/auth/authenticator"
	"car_wash/infra/mux/controller"
	"github.com/gorilla/mux"
	"net/http"
)

func InitRouter(controller *controller.Controller, authHandler *authenticator.Authenticator) *mux.Router {
	router := mux.NewRouter()

	authenticationRoutes := router.PathPrefix("/auth").Subrouter()
	authenticationRoutes.HandleFunc("/register", controller.RegisterNewOwner).Methods(http.MethodPost)
	authenticationRoutes.HandleFunc("/cache", controller.AddCredsToJar).Methods(http.MethodPost)
	authenticationRoutes.Use(authHandler.BearerAuth)

	pathSubrouter := router.PathPrefix("/carwash").Subrouter()
	pathSubrouter.HandleFunc("/wss", controller.UpgradeWss).Methods(http.MethodGet)
	pathSubrouter.HandleFunc("/register", controller.RegisterCarWash).Methods(http.MethodPost)
	pathSubrouter.HandleFunc("/wash", controller.RegisterWash).Methods(http.MethodPost)
	pathSubrouter.Use(authHandler.APIAuth)

	router.HandleFunc("/auth/check/{hash}", controller.CheckHash).Methods(http.MethodGet)

	return router
}
