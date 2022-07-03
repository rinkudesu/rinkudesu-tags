package Controllers

import (
	"net/http"
)

//todo: look into some sort of DI

type TagsController struct {
}

func (controller *TagsController) HandleProducts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(420)
	return
}
