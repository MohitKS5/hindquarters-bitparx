package util

import (
	"net/http"
	"time"
	"math/rand"
	"encoding/json"
	"fmt"
	"log"
)

// JSONResponse represents an HTTP response which contains a JSON body.
type JSONResponse struct {
	// HTTP status code.
	Code int
	// JSON represents the JSON that should be serialized and sent to the client
	JSON interface{}
	// Headers represent any headers that should be sent to the client
	Headers map[string]string
}

func (res JSONResponse) Encode(w *http.ResponseWriter) {
	err, ok := res.JSON.(ParxError)
	if ok {
		http.Error(*w, err.Err, res.Code)
	} else {
		encerr := json.NewEncoder(*w).Encode(res)
		if encerr!=nil{
			log.Println(encerr)
		}
	}
}

// Error represents the "standard error response"
type ParxError struct {
	ErrCode string `json:"errcode"`
	Err     string `json:"error"`
}

func (e *ParxError) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrCode, e.Err)
}

func LogThenError(err error, where ...string) error {
	log.Println(where, ": ",err)
	return err
}

// MessageResponse returns a JSONResponse with a 'message' key containing the given text.
func MessageResponse(code int, msg string) JSONResponse {
	return JSONResponse{
		Code: code,
		JSON: struct {
			Message string `json:"message"`
		}{msg},
	}
}

// ErrorResponse returns an HTTP 500 JSONResponse with the stringified form of the given error.
func ErrorResponse(err error) JSONResponse {
	return MessageResponse(500, err.Error())
}

// SetCORSHeaders sets unrestricted origin Access-Control headers on the response writer
func SetCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

const alphanumerics = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// RandomString generates a pseudo-random string of length n.
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphanumerics[rand.Int63()%int64(len(alphanumerics))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}