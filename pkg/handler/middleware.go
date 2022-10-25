package handler

import (
	"backend/pkg/utils"
	"context"
	"errors"
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
		tokenString := request.Header["Authorization"]
		if len(tokenString) == 0 {
			NewErrorResponse(writer, request, errors.New("empty token"), 404)
			return
		} else {
			token := strings.Split(request.Header["Authorization"][0], " ")
			id, err := utils.ParseToken(token[1])

			if err != nil {
				NewErrorResponse(writer, request, err, 400)
				return
			}

			var UserID = "UserID"
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), UserID, id)))

			//TODO if user exist put them to request context
		}
	})
}
