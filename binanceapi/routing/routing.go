package routing

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

const (
	BASE_URL = "https://api.binance.com"
)

func Setup(router *mux.Router) {
	s := router.PathPrefix("/trade").Subrouter()
	s.HandleFunc("/me", test)
	s.HandleFunc("/account", serve(getAccount))
	s.HandleFunc("/depth", serve(Depth))
	s.HandleFunc("/exchange", serve(GetExchangeInfo))
	router.Use(loggerMiddleware)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println("route: ", request.URL.Path)
		next.ServeHTTP(writer, request)
	})
}

func test(w http.ResponseWriter, r *http.Request) {

}
