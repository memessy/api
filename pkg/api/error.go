package api

import (
	"encoding/json"
	"net/http"
)

type ErrResponse struct {
	HttpStatusCode int         `json:"-"`
	Message        string      `json:"message"`
	Detail         interface{} `json:"detail"`
}

func (r ErrResponse) RenderJson(w http.ResponseWriter) {
	w.WriteHeader(r.HttpStatusCode)
	w.Header().Add("content-type", "application/json")
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		panic("err response error")
	}
}

func InternalError(w http.ResponseWriter) {
	ErrResponse{
		HttpStatusCode: http.StatusInternalServerError,
		Message:        "something went wrong",
	}.RenderJson(w)
}

func ValidationError(w http.ResponseWriter, message string, detail interface{}) {
	ErrResponse{
		HttpStatusCode: http.StatusUnprocessableEntity,
		Message:        message,
		Detail:         detail,
	}.RenderJson(w)
}
