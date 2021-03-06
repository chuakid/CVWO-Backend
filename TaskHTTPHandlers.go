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
	task, err := extractTaskAndCheckAccess(r)
	if err != nil {
		log.Println(err)
		if err.Error() == "not auth" {
			http.Error(w, "Not authorised", 401)
		} else {
			http.Error(w, "Error editing task", 400)
		}
		return
	}

	descStruct := struct {
		Description string `json:"description"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&descStruct)
	if err != nil {
		log.Println("Error editing task", err)
		http.Error(w, "Error editing task", 400)
		return
	}

	err = task.EditTask(descStruct.Description)
	if err != nil {
		log.Println("Error editing task", err)
		http.Error(w, "Error editing task", 400)
		return
	}
	w.Write([]byte("Task edited"))
}

func setTaskCompletionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Set task completion handler hit")
	task, err := extractTaskAndCheckAccess(r)
	if err != nil {
		log.Println(err)
		if err.Error() == "not auth" {
			http.Error(w, "Not authorised", 401)
		} else {
			http.Error(w, "Error editing task", 400)
		}
		return
	}
	completedStruct := struct {
		Completed bool `json:"completed"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&completedStruct)
	if err != nil {
		log.Println("Error editing task", err)
		http.Error(w, "Error editing task", 400)
		return
	}

	task.Completed = completedStruct.Completed

	err = task.SetTaskCompletion()

	if err != nil {
		log.Println("Error editing task", err)
		http.Error(w, "Error editing task", 400)
		return
	}
	w.Write([]byte("Task edited"))
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get tasks endpoint hit")
	userid := r.Context().Value("userid")
	if userid, ok := userid.(string); ok { //Type assertion
		useridInt, err := strconv.Atoi(userid)
		if err != nil {
			log.Println("Error getting tasks:", err)
			http.Error(w, "Error getting tasks", 400)
			return
		}
		user := models.User{ID: useridInt}
		tasks, err := user.GetTasks()
		if err != nil {
			log.Println("Error getting tasks:", err)
			http.Error(w, "Error getting tasks", 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}

}

func changeColorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Change color endpoint hit")
	task, err := extractTaskAndCheckAccess(r)
	if err != nil {
		log.Println(err)
		if err.Error() == "not auth" {
			http.Error(w, "Not authorised", 401)
		} else {
			http.Error(w, "Error changing color of task", 400)
		}
		return
	}

	colorStruct := struct {
		Color int `json:"color"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&colorStruct)
	if err != nil {
		log.Println("Error changing color of task", err)
		http.Error(w, "Error changing color of task", 400)
		return
	}

	err = task.ChangeColor(colorStruct.Color)
	if err != nil {
		log.Println("Error changing color of task", err)
		http.Error(w, "Error changing color of task", 400)
		return
	}
	w.Write([]byte("Color changed"))
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
