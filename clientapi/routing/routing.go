package routing

import (
	"github.com/gorilla/mux"
	"net/http"
)

// routes comfigured here

func Setup(router *mux.Router)  {
	router.HandleFunc("/welcome", SayWelcome).Methods(http.MethodGet)
}