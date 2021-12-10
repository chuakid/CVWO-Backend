package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
)

//Router for protected routes
func ProtectedRoutes() *chi.Mux {
	protectedR := chi.NewRouter()
	protectedR.Use(loggedInOnly)
	protectedR.Route("/project", func(r chi.Router) {
		r.Get("/{projectId}", getProjectHandler)
		r.Post("/", createProjectHandler)
	})
	protectedR.Get("/projects", getProjectsHandler)
	protectedR.Route("/task", func(r chi.Router) {
		r.Post("/", createTaskHandler)
	})
	return protectedR
}

func loggedInOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Not logged in", 401)
			return
		}
		tokensplit := strings.Split(authHeader, "Bearer ")
		if len(tokensplit) < 2 {
			http.Error(w, "Not logged in", 401)
			return
		}
		tokenstring := tokensplit[1]
		token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is hmac:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, nil
			}
			return []byte(key), nil
		})
		if err != nil {
			log.Print(err)
			http.Error(w, "Invalid token", 401)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "userid", claims["sub"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			log.Print(err)
			return
		}
	})
}
