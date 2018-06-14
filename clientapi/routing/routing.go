package routing

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/bitparx/clientapi/auth/storage/accounts"
	"github.com/bitparx/clientapi/auth/storage/devices"
	"github.com/bitparx/clientapi/auth"
	"github.com/bitparx/common/jsonerror"
	"github.com/bitparx/clientapi/auth/storage/levels"
	"database/sql"
	"github.com/gorilla/context"
	"github.com/bitparx/clientapi/Integration_apis/proxy_handles"
)

// config route
type routerConfig struct {
	// map for endpoints with auth bypass = true
	routeAuth map[string]bool
}

// routes comfigured here usign gorilla/mux
// encaspulated handlers are used to provide with database pointers

func Setup(router *mux.Router, accountDB *accounts.Database, deviceDB *devices.Database, levelDB *levels.Database) {
	route := routerConfig{map[string]bool{
		"/login":    true,
		"/register": true,
		"/welcome":  true,
	}}
	router.HandleFunc("/welcome", SayWelcome).Methods(http.MethodGet)
	router.HandleFunc("/login", LoginHandler(accountDB, deviceDB, levelDB)).Methods(http.MethodPost)
	router.HandleFunc("/register", RegisterHandler(accountDB, deviceDB, levelDB)).Methods(http.MethodPost)
	router.HandleFunc("/registration", RegistrationHandler(levelDB)).Methods(http.MethodPut, http.MethodDelete, http.MethodGet)

	// routes with auth = true
	router.HandleFunc("/logout", LogoutHandler(deviceDB)).Methods(http.MethodPost)
	router.HandleFunc("/levels", RouteLevelsHandler(levelDB, sql.NullBool{false, true}, GetAllAccountLevels)).Methods(http.MethodPost)
	router.HandleFunc("/levels/{levelname}", RouteLevelsHandler(levelDB, sql.NullBool{}, RequestLevelByLocalpart)).Methods(http.MethodPost)
	router.HandleFunc("/levels/{levelname}/{localpart}", RouteLevelsHandler(levelDB, sql.NullBool{true, true}, SetLevelByLocalpart)).Methods(http.MethodPut)
	router.HandleFunc("/levels/{levelname}/{localpart}", RouteLevelsHandler(levelDB, sql.NullBool{false, true}, SetLevelByLocalpart)).Methods(http.MethodDelete)
	router.HandleFunc("/accounts", RouteHandlerAccounts(accountDB)).Methods(http.MethodPost)
	router.HandleFunc("/devices", RouteHandlerDevices(deviceDB)).Methods(http.MethodPost)

	//binance api paths
	binanceProxy  := proxy_handles.NewProxy("localhost:8080")
	router.HandleFunc("/trade",binanceProxy.Handle)

	router.Use(route.authMiddleware(deviceDB))
}

// middleware for auth=true endpoints, uses routeConfig and lets bypass the true mapped routes
// adds localpart of client to the context with "localpart" key
func (route routerConfig) authMiddleware(deviceDB *devices.Database) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if route.routeAuth[request.URL.Path] {
				next.ServeHTTP(writer, request)
			} else {
				dev, err := auth.VerifyAccessToken(request, deviceDB)
				if err != nil {
					myerr, ok := err.JSON.(*jsonerror.ParxError)
					if ok {
						http.Error(writer, myerr.Err, err.Code)
					}
					return
				} else {
					context.Set(request, "localpart", dev.UserID)
					next.ServeHTTP(writer, request)
				}
			}
		})
	}
}
