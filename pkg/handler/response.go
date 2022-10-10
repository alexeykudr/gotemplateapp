package handler

import (
	"encoding/json"
	"net/http"
)

// NewErrorResponse Send error and status as json response
func NewErrorResponse(w http.ResponseWriter, r *http.Request, err error, status int) {
	res, _ := json.Marshal(err)
	w.WriteHeader(status)
	w.Write(res)
}
