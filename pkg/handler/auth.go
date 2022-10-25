package handler

import (
	"backend/domain"
	"context"
	"encoding/json"
	"net/http"
)

type signInInput struct {
	Email    string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		NewErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}

	err = user.Validate()
	if err != nil {
		NewErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateUser(context.Background(), user)
	if err != nil {
		NewErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}

	resp := make(map[string]int)
	resp["id"] = id
	jsonResp, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	v := r.Context().Value("UserID")
	r_id := r.Context().Value("RequestIDKey")
	resp := make(map[string]string)
	resp[v.(string)] = r_id.(string)
	jsonResp, _ := json.Marshal(resp)
	NewOkResponse(w, r, jsonResp)

}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var inputUser signInInput
	err := json.NewDecoder(r.Body).Decode(&inputUser)
	if err != nil {
		NewErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	token, err := h.Service.Login(inputUser.Email, inputUser.Password)
	if err != nil {
		NewErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}
	resp := make(map[string]string)
	resp["token"] = token
	jsonResp, _ := json.Marshal(resp)
	NewOkResponse(w, r, jsonResp)
}
