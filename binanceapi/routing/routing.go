package routing

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"github.com/bitparx/common/config"
)

const (
	BASE_URL = config.BINANCE_REST_URL
)

func Setup(router *mux.Router) {
	s := router.PathPrefix("/trade").Subrouter()
	s.HandleFunc("/welcome", test)
	s.HandleFunc("/account", serve(getAccount))
	s.HandleFunc("/order/all", serve(ListAllOrders))
	s.HandleFunc("/depth", serve(Depth))
	s.HandleFunc("/exchange", serve(GetExchangeInfo))
	s.HandleFunc("/mytrade", serve(getMyTrades))
	s.HandleFunc("/order/open", serve(ListOpenOrders))
	s.HandleFunc("/order", serve(GetOrderById))
	s.HandleFunc("/reports/aggregate", serve(GetAggregateTrades))
	s.HandleFunc("/reports/recent", serve(GetRecentTrades))
	router.Use(loggerMiddleware)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println("route: ", request.URL.Path)
		next.ServeHTTP(writer, request)
	})
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello there"))
}
