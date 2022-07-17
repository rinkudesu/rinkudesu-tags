package Controllers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Ok(w http.ResponseWriter) {
	log.Debug("Responding with OK")
	w.WriteHeader(200)
}

func BadRequest(w http.ResponseWriter) {
	log.Debug("Responding with BAD REQUEST")
	w.WriteHeader(400)
}

func InternalServerError(w http.ResponseWriter) {
	log.Debug("Responding with INTERNAL ERROR")
	w.WriteHeader(500)
}

func NotFound(w http.ResponseWriter) {
	log.Debug("Responding with NOT FOUND")
	w.WriteHeader(404)
}
