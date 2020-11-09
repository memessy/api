package api

import (
	"encoding/json"
	"net/http"
)

func RenderJson(object interface{}, w http.ResponseWriter) error {
	w.Header().Add("content-type", "application/json")
	return json.NewEncoder(w).Encode(object)
}
