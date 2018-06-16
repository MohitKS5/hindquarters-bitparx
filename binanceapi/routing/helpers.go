package routing

import (
	"github.com/bitparx/util"
	"reflect"
	"github.com/bitparx/common/jsonerror"
	"net/http"
	"net/url"
	"io/ioutil"
	"errors"
)

type bnbreq struct {
	jsonres util.JSONResponse
}

func DialBnb(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, util.LogThenError(err, "request")
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, util.LogThenError(errors.New(bodyString),"error: ")
	}
	return resp, err
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
