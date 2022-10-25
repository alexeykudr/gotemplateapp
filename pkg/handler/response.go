package handler

import (
	"encoding/json"
	"net/http"
)

// NewErrorResponse Send error and status as json response
func NewErrorResponse(w http.ResponseWriter, r *http.Request, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["Error"] = err.Error()
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(status)
	w.Write(jsonResp)
}

func NewOkResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]interface{})
	resp["Ok"] = data
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
