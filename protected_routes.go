package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chuakid/cvwo-backend/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
)

//Router for protected routes
func ProtectedRoutes() *chi.Mux {
	protectedR := chi.NewRouter()
	protectedR.Use(loggedInOnly)
	protectedR.Route("/project", func(r chi.Router) {
		r.Get("/{projectId}", getProjectHandler)
		r.Post("/", uploadProject)
	})
	protectedR.Get("/projects", getProjects)
	protectedR.Route("/task", func(r chi.Router) {
		r.Post("/", createTaskHandler)
	})
	return protectedR
}

func getProjects(w http.ResponseWriter, r *http.Request) {
	log.Print("Get Projects Endpoint Hit")
	userid := r.Context().Value("userid")
	if userid, ok := userid.(string); ok { //Type assertion
		id, err := strconv.Atoi(userid)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error getting projects", 400)
			return
		}
		user := models.User{ID: id}
		projects, err := user.GetProjects()
		if err != nil {
			http.Error(w, "Error getting projects", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(projects)
	} else {
		return
	}
}

func uploadProject(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value("userid")
	if userid, ok := userid.(string); ok { //Type assertion
		var project models.Project
		err := json.NewDecoder(r.Body).Decode(&project)
		if err != nil {
			log.Print("Error uploading project:", err)
			http.Error(w, "Error uploading project", 400)
			return
		}
		_, err = project.CreateProject(userid)
		if err != nil {
			log.Print("Error uploading project:", err)
			http.Error(w, "Error uploading project", 400)
			return
		}
		w.Header().Set("Cotent-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":   fmt.Sprint(project.ID),
			"name": project.Name,
		})

	} else {
		log.Print("Error uploaded project:")
		return
	}
}

func getProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Get Tasks Endpoint Hit")
	userid := r.Context().Value("userid")
	if userid, ok := userid.(string); ok { //Type assertion
		projectid := chi.URLParam(r, "projectId")
		projectidInt, err := strconv.Atoi(projectid)
		if err != nil {
			log.Println("Error getting project:", err)
			http.Error(w, "Error getting project", 400)
			return
		}
		project := models.Project{ID: projectidInt}
		//check if user is allowed to access project
		projectAccess, err := checkProjectAccess(project, userid)
		if err != nil {
			log.Println("Error getting project:", err)
			http.Error(w, "Error getting project", 400)
			return
		}
		if !projectAccess {
			http.Error(w, "Not authorized to access project", 401)
			return
		}
		err = project.GetProject()
		if err != nil {
			log.Println("Error getting project:", err)
			http.Error(w, "Error getting project", 400)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(project)
	} else {
		return
	}
}

func checkProjectAccess(project models.Project, userid string) (bool, error) {
	projectUsers, err := project.GetUsers()
	if err != nil {
		return false, err
	}
	for _, user := range projectUsers {
		if fmt.Sprint(user.ID) == userid {
			return true, nil
		}
	}
	log.Println("Not authorized to access project")
	return false, nil

}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
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
