package routing

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/bitparx/clientapi/auth/storage/accounts"
	"github.com/bitparx/clientapi/auth/storage/devices"
	"github.com/bitparx/clientapi/auth"
	"github.com/bitparx/common/jsonerror"
	"github.com/bitparx/clientapi/auth/storage/levels"
)

// config route
type routerConfig struct {
	// map for endpoints with auth = true
	routeAuth map[string]bool
}

// routes comfigured here usign gorilla/mux
// encaspulated handlers are used to provide with database pointers

func Setup(router *mux.Router, accountDB *accounts.Database, deviceDB *devices.Database, levelDB *levels.Database) {
	route := routerConfig{map[string]bool{
		"/logout": true,
		"/levels": false,
	}}
	router.HandleFunc("/welcome", SayWelcome).Methods(http.MethodGet)
	router.HandleFunc("/login", LoginHandler(accountDB, deviceDB, levelDB)).Methods(http.MethodPost)
	router.HandleFunc("/register", RegisterHandler(accountDB, deviceDB, levelDB)).Methods(http.MethodPost)

	// routes with auth = true
	router.HandleFunc("/logout", LogoutHandler(deviceDB)).Methods(http.MethodPost)
	router.HandleFunc("/levels", RouteLevelsHandler(levelDB, false, GetAllAccountLevels)).Methods(http.MethodPost)
	router.HandleFunc("/levels/{levelname}/{localpart}", RouteLevelsHandler(levelDB, true, UpdateLevelByLocalpart)).Methods(http.MethodPut)
	router.HandleFunc("/levels/{levelname}/{localpart}", RouteLevelsHandler(levelDB, false,  UpdateLevelByLocalpart)).Methods(http.MethodDelete)
	router.Use(route.authMiddleware(deviceDB))
}

// middleware for auth=true endpoints
func (route routerConfig) authMiddleware(deviceDB *devices.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if route.routeAuth[request.URL.Path] {
				_, err := auth.VerifyAccessToken(request, deviceDB)
				if err != nil {
					myerr, ok := err.JSON.(*jsonerror.ParxError)
					if ok {
						http.Error(writer, myerr.Err, err.Code)
					}
					return
				} else {
					next.ServeHTTP(writer, request)
				}
			} else {
				next.ServeHTTP(writer, request)
			}
		})
	}
}
