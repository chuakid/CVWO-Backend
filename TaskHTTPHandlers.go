package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/chuakid/cvwo-backend/models"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
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

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := extractTaskAndCheckAccess(r)
	if err != nil {
		log.Println(err)
		if err.Error() == "not auth" {
			http.Error(w, "Not authorised", 401)
		} else {
			http.Error(w, "Error deleting task", 400)
		}
		return
	}

	err = task.DeleteTask()
	if err != nil {
		log.Println("Error deleting task", err)
		http.Error(w, "Error deleting task", 400)
		return
	}
	w.Write([]byte("Task deleted"))

}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func extractTaskAndCheckAccess(r *http.Request) (*models.Task, error) {
	userid := r.Context().Value("userid")
	if userid, ok := userid.(string); ok { //Type assertion
		taskId := chi.URLParam(r, "taskId")
		taskIdInt, err := strconv.Atoi(taskId)
		if err != nil {
			return nil, err
		}

		task := models.Task{
			ID: taskIdInt,
		}
		err = task.GetTask()
		if err != nil {
			return nil, err
		}

		project := models.Project{ID: task.ProjectID}
		projectAccess, err := checkProjectAccess(project, userid)
		if err != nil {
			return nil, err
		}
		if !projectAccess {
			return nil, errors.New("not auth")
		}
		return &task, nil
	} else {
		return nil, errors.New("Error extracting taskid")
	}
}
