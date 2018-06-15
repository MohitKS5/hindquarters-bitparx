package routing

import (
	"github.com/bitparx/util"
	"reflect"
	"github.com/bitparx/common/jsonerror"
	"net/http"
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
func Invoke(fn interface{}, args ...string) ([]reflect.Value) {
	v := reflect.ValueOf(fn)
	rargs := make([]reflect.Value, len(args))
	for i, a := range args {
		rargs[i] = reflect.ValueOf(a)
	}
	res := v.Call(rargs)
	return res
}

// serve
func serve(fn interface{}, args ...string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		r := wrap(Invoke(fn, args...))
		r.jsonres.Encode(&writer)
	}
}

