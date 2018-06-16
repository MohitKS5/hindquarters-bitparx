package routing

import (
	"github.com/bitparx/util"
	"reflect"
	"github.com/bitparx/common/jsonerror"
	"net/http"
	"net/url"
)

type bnbreq struct {
	jsonres util.JSONResponse
}

// make json response
func wrap(args []reflect.Value) (r bnbreq) {
	r = bnbreq{}
	if ok := args[1]; !ok.IsNil() {
		r.jsonres = util.JSONResponse{
			Code: 500,
			JSON: jsonerror.InternalServerError(),
		}
		return
	}
	r.jsonres = util.JSONResponse{
		Code: 200,
		JSON: args[0].Interface(),
	}
	return
}

// invoke an anonymous function
func Invoke(fn interface{}, args url.Values) ([]reflect.Value) {
	v := reflect.ValueOf(fn)
	rargs := []reflect.Value{reflect.ValueOf(args)}
	res := v.Call(rargs)
	return res
}

// serve
func serve(fn interface{}) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		r := wrap(Invoke(fn, request.URL.Query()))
		r.jsonres.Encode(&writer)
	}
}
