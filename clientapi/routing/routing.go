package routing

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/bitparx/clientapi/auth/storage/accounts"
	"github.com/bitparx/clientapi/auth/storage/devices"
	"encoding/json"
	"github.com/bitparx/clientapi/auth"
)

// config route
type routerConfig struct {
	// map for endpoints with auth = true
	routeAuth map[string]bool
}

// routes comfigured here usign gorilla/mux
// encaspulated handlers are used to provide with database pointers

func Setup(router *mux.Router, accountDB *accounts.Database, deviceDB *devices.Database) {
	route:= routerConfig{map[string]bool{
		"/logout": true,
	}}
	router.HandleFunc("/welcome", SayWelcome).Methods(http.MethodGet)
	router.HandleFunc("/login", LoginHandler(accountDB, deviceDB)).Methods(http.MethodPost)
	router.HandleFunc("/register", RegisterHandler(accountDB, deviceDB)).Methods(http.MethodPost)

	// routes with auth = true
	router.HandleFunc("/logout", LogoutHandler(deviceDB)).Methods(http.MethodPost)
	router.Use(route.authMiddleware(deviceDB))
}


// middleware for auth=true endpoints
func (route routerConfig) authMiddleware(deviceDB *devices.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if route.routeAuth[request.URL.Path] {
				_, err := auth.VerifyAccessToken(request, deviceDB)
				if err != nil {
					json.NewEncoder(writer).Encode(err)
					return
				}
			}
			next.ServeHTTP(writer, request)
		})
	}
}
