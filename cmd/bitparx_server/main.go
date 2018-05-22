package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/gorilla/handlers"
	"github.com/bitparx/clientapi/routing"
	//"github.com/bitparx/common/storage"
)

func main()  {
	fmt.Println("Starting server at http://localhost:12345...")

	// setting up router
	router := mux.NewRouter()
	routing.Setup(router)

	//databse setup
	//db := storage.PostgresConnect()

	// starting server while setting cors for angular
	log.Fatal(http.ListenAndServe(":12345", handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD"}),
		handlers.AllowedOrigins([]string{"*"}))(router)))
}