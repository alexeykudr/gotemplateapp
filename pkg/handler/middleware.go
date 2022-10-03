package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func (h *Handler) RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.NewUUID()
		var requestID = id
		//https://medium.com/golangspec/globally-unique-key-for-context-value-in-golang-62026854b48f
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestID, id)))
	})
}
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		a := "admin"
		b := "123"
		if ok {
			if username == a && password == b {
				next.ServeHTTP(writer, request)
			} else {
				writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			}
		}
	})
}

func (h *Handler) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		token := strings.Split(request.Header["Authorization"][0], " ")
		fmt.Println(token[1])
		id, err := h.Service.ParseToken(token[1])
		fmt.Println(id)
		fmt.Println(err)
	})
}
