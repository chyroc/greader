package greader_api

import "net/http"

func WrapStdHttpRequest(request *http.Request) HttpReader {
	return &httpReader{request}
}

type httpReader struct {
	request *http.Request
}

func (r *httpReader) Form(key string) ([]string, bool) {
	val, ok := r.request.Form[key]
	return val, ok
}

func (r *httpReader) Header(key string) ([]string, bool) {
	val, ok := r.request.Header[key]
	return val, ok
}

func (r *httpReader) FormString(key string) string {
	val, ok := r.request.Form[key]
	if !ok {
		return ""
	}
	return val[0]
}

func (r *httpReader) QueryString(key string) string {
	val, ok := r.request.URL.Query()[key]
	if !ok {
		return ""
	}
	return val[0]
}
