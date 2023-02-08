package greader_api

import "net/http"

func WrapStdHttpRequest(request *http.Request) HttpReader {
	return &stdHttpReader{request}
}

type stdHttpReader struct {
	request *http.Request
}

func (r *stdHttpReader) HeaderString(key string) string {
	val, _ := r.request.Header[key]
	if len(val) > 0 {
		return val[0]
	}
	return ""
}

func (r *stdHttpReader) FormList(key string) []string {
	val, _ := r.request.Form[key]
	return val
}

func (r *stdHttpReader) FormString(key string) string {
	val, ok := r.request.Form[key]
	if !ok {
		return ""
	}
	return val[0]
}

func (r *stdHttpReader) QueryString(key string) string {
	val, ok := r.request.URL.Query()[key]
	if !ok {
		return ""
	}
	return val[0]
}
