package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/bitparx/binanceapi/routing"
)

func main() {
	fmt.Println("Starting server at http://localhost:8080")

	// setting up router
	router := mux.NewRouter()
	routing.Setup(router)

	// starting server while setting cors for angular
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE"}),
		handlers.AllowedOrigins([]string{"http://localhost:12345"}))(router)))
}
