package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"log"

	"github.com/chuakid/cvwo-backend/db"
	"github.com/chuakid/cvwo-backend/models"
	"github.com/go-chi/chi/v5"
)

func main() {
	//Set up database
	err := db.InitDatabase()
	if err != nil {
		log.Fatalln("could not create database", err)
	}
	db.DB.AutoMigrate(&models.User{})
	db.DB.AutoMigrate(&models.Project{})

	//Set up logger
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(file)

	//Set up router and routes
	r := chi.NewRouter()
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Login endpoint hit")
	})

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Register endpoint hit")

		var u models.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Println("Error while decoding: ", err)
			http.Error(w, http.StatusText(400), 400)
			return
		}

		err = register(u.Username, u.Password)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: users.username" { //username taken
				http.Error(w, "Username taken", 400)
			} else {
				log.Println("Error: ", err)
				http.Error(w, "Error", 400)
			}
			return
		}
		w.Write([]byte("Success"))
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/project", getProject)
	})

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8000"
	}
	fmt.Println("Listening on port:", PORT)
	http.ListenAndServe(":"+PORT, r)
}

func getProject(w http.ResponseWriter, r *http.Request) {

}

func register(username string, password string) error {
	user := models.User{
		Username: username,
		Password: password,
	}
	id, err := user.CreateUser()
	if err != nil {
		return err
	}
	log.Print("User created:", id)
	return nil
}
