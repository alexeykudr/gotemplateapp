package handler

import (
	"backend"
	"context"
	"encoding/json"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *signInInput) Validate() error {
	return validation.ValidateStruct(&s,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&s.Username, validation.Required, validation.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&s.Password, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
	)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user backend.User
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

func (h *Handler) SecretInfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "secret info")
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var input signInInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		NewErrorResponse(w, r, err, http.StatusBadRequest)
		return
	}

	token, err := h.Service.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	resp := make(map[string]string)
	resp["token"] = token
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
