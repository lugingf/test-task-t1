package access_logger

import (
	"net/http"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status       int
	ResponseBody string
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *StatusRecorder) Write(body []byte) (int, error) {
	r.ResponseBody = string(body)
	return r.ResponseWriter.Write(body)
}

