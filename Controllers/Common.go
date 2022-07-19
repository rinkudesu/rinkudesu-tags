package Controllers

import (
	json2 "encoding/json"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"io"
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

func MethodNotAllowed(w http.ResponseWriter) {
	log.Debug("Responding with METHOD NOT ALLOWED")
	w.WriteHeader(405)
}

func ReadBody(body io.ReadCloser) ([]byte, error) {
	array, err := io.ReadAll(body)
	if err != nil {
		log.Warningf("Failed to read from body %s", err.Error())
	}
	return array, err
}

func CloseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		log.Warningf("Failed to close request body: %s", err.Error())
	}
}
func WriteJsonResponse(w http.ResponseWriter, code int, object interface{}) {
	json, jsonErr := json2.Marshal(object)
	if jsonErr != nil {
		log.Warningf("Failed to serialise to json: %s", jsonErr.Error())
		InternalServerError(w)
		return
	}
	w.WriteHeader(code)
	_, err := w.Write(json)
	if err != nil {
		log.Warningf("Failed to write response: %s", err.Error())
		InternalServerError(w)
		return
	}
}

func ParseUuid(id string) (uuid.UUID, error) {
	result, err := uuid.FromString(id)
	if err != nil {
		log.Infof("Unable to parse '%s' as uuid", id)
		return uuid.Nil, err
	}
	return result, nil
}
