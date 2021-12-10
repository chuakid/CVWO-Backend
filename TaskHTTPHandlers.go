package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chuakid/cvwo-backend/models"
)

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Create tasks endpoint hit")
	userid := r.Context().Value("userid")
	if userid, ok := userid.(string); ok { //Type assertion
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Println("Error creating task:", err)
			http.Error(w, "Error creating task", 400)
			return
		}
		projectAccess, err := checkProjectAccess(models.Project{ID: task.ProjectID}, userid)
		if err != nil {
			log.Println("Error creating task:", err)
			http.Error(w, "Error creating task", 400)
			return
		}
		if !projectAccess {
			http.Error(w, "Not authorized to access project", 401)
			return
		}

		err = task.CreateTask()
		if err != nil {
			log.Println("Error creating task:", err)
			http.Error(w, "Error creating task", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}

}
