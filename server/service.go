package server

import (
	"net/http"
)

type Service interface {
	GetTime(http.ResponseWriter, *http.Request)

	PostSignup(http.ResponseWriter, *http.Request)

	PostLogin(http.ResponseWriter, *http.Request)

	PostWebhook(http.ResponseWriter, *http.Request)
}
