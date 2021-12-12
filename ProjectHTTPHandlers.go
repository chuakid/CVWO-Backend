package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/chuakid/cvwo-backend/models"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
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
	log.Print("Get Project Endpoint Hit")
	project, err := extractProjectIdAndCheckAccess(r)
	if err != nil {
		log.Println(err)
		if err.Error() == "not auth" {
			http.Error(w, "Not authorised", 401)
		} else {
			http.Error(w, "Error getting project", 400)
		}
		return
	}
	apiproject, err := project.GetProjectDetails()
	if err != nil {
		log.Println("Error getting project:", err)
		http.Error(w, "Error getting project", 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiproject)

}

func deleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Delete Project Endpoint Hit")
	project, err := extractProjectIdAndCheckAccess(r)
	if err != nil {
		log.Println(err)
		if err.Error() == "not auth" {
			http.Error(w, "Not authorised", 401)
		} else {
			http.Error(w, "Error deleting project", 400)
		}
		return
	}
	err = project.DeleteProject()
	if err != nil {
		log.Println("Error deleting project:", err)
		http.Error(w, "Error deleting project", 400)
		return
	}
	w.Write([]byte("Deletion success"))

}

func renameProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Rename Project Endpoint Hit")
	project, err := extractProjectIdAndCheckAccess(r)
	if err != nil {
		log.Println(err)
		if err.Error() == "not auth" {
			http.Error(w, "Not authorised", 401)
		} else {
			http.Error(w, "Error renaming project", 400)
		}
		return
	}
	namestruct := struct {
		Name string `json:"name"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&namestruct)
	if err != nil {
		log.Println("Error renaming project:", err)
		http.Error(w, "Error renaming project", 400)
	}
	err = project.RenameProject(namestruct.Name)
	if err != nil {
		log.Println("Error renaming project:", err)
		http.Error(w, "Error renaming project", 400)
	}
	w.Write([]byte("Rename success"))

}

func addProjectUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Add Project User Endpoint Hit")
	project, err := extractProjectIdAndCheckAccess(r)
	if err != nil {
		log.Println(err)
		if err.Error() == "not auth" {
			http.Error(w, "Not authorised", 401)
		} else {
			http.Error(w, "Error adding user to project", 400)
		}
		return
	}

	//Check if user exists
	username := chi.URLParam(r, "username")
	user := models.User{Username: username}
	if !user.UserExists() {
		log.Println("User does not exist")
		http.Error(w, "User does not exist", 400)
		return
	}

	err = project.AddUser(&user)
	if err != nil {
		log.Println("Error adding user to project:", err)
		http.Error(w, "Error adding user to project", 400)
	}
	w.Write([]byte("User added"))
}

func extractProjectIdAndCheckAccess(r *http.Request) (*models.Project, error) {
	userid := r.Context().Value("userid")
	if userid, ok := userid.(string); ok { //Type assertion
		projectid := chi.URLParam(r, "projectId")
		projectidInt, err := strconv.Atoi(projectid)
		if err != nil {
			return nil, err
		}
		project := models.Project{ID: projectidInt}
		projectAccess, err := checkProjectAccess(project, userid)
		if err != nil {
			return nil, err
		}
		if !projectAccess {
			return nil, errors.New("not auth")
		}
		return &project, nil
	} else {
		return nil, errors.New("Error extracting userid")
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
