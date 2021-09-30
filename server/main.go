package server

import (
	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, s Service) {
	router.HandleFunc("/", s.GetTime).Methods("GET")

	router.HandleFunc("/login", s.PostLogin).Methods("POST")

	router.HandleFunc("/signup", s.PostSignup).Methods("POST")

	router.HandleFunc("/webhook", s.PostWebhook).Methods("POST")

}
