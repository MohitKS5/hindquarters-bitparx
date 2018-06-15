package routing

import (
	"net/http"
	"fmt"
	"encoding/json"
)

//test function for server

func SayWelcome(w http.ResponseWriter, req *http.Request) {
	fmt.Println("path", req.URL.Path)
	json.NewEncoder(w).Encode("hello bitparx")
}
