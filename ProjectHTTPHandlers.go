package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/chuakid/cvwo-backend/models"
	"github.com/go-chi/chi/v5"
)

func getProjectsHandler(w http.ResponseWriter, r *http.Request) {
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

func createProjectHandler(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value("userid")
	if userid, ok := userid.(string); ok { //Type assertion
		var project models.Project
		err := json.NewDecoder(r.Body).Decode(&project)
		if err != nil {
			log.Print("Error uploading project:", err)
			http.Error(w, "Error uploading project", 400)
			return
		}
		err = project.CreateProject(userid)
		if err != nil {
			log.Print("Error uploading project:", err)
			http.Error(w, "Error uploading project", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
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
