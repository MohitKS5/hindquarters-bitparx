package utils

type ReqParams map[string]string

func MergeMaps(params ...ReqParams) (query map[string]string) {
	for k := range params {
		for v := range params[k] {
			query[v] = params[k][v]
		}
	}
	return
}
