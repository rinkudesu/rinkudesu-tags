package Controllers

import "net/http"

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(400)
}

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(500)
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(404)
}
