package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/rs/cors"
	"github.com/gorilla/mux"
)

type Task struct {
	ID     		string 	`json:"id,omitempty"`
	Title 		string 	`json:"title,omitempty"`
	Status		string 	`json:"status,omitempty"`
}

var tasks []Task

func GetTodosEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func GetTaskEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	for _, item := range tasks {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(struct {

	}{})
}

func CreateTaskEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var task Task
	_ = json.NewDecoder(req.Body).Decode(&task)
	task.ID = params["id"]
	tasks  = append(tasks, task)
	json.NewEncoder(w).Encode(tasks)
}

func DeleteTaskEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	for index, item := range tasks {
		for item.ID == params["id"]{
			tasks = append(tasks[:index], tasks[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(tasks)
}

func UpdateTaskEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var task Task
	_ = json.NewDecoder(req.Body).Decode(&task)

	for index, item := range tasks {
		if item.ID == params["id"]{
			tasks[index].Status = task.Status
			break
		}
	}

	json.NewEncoder(w).Encode(tasks)
}

func main() {
	router := mux.NewRouter()

	tasks = append(tasks, Task{ID: "1", Title: "DB Schema", Status: "true"})
	tasks = append(tasks, Task{ID: "2", Title: "Markup", Status: "false"})

	router.HandleFunc("/todos", GetTodosEndpoint).Methods("GET")
	router.HandleFunc("/todos/{id}", GetTaskEndpoint).Methods("GET")
	router.HandleFunc("/todos/{id}", CreateTaskEndpoint).Methods("POST")
	router.HandleFunc("/todos/{id}", DeleteTaskEndpoint).Methods("DELETE")
	router.HandleFunc("/todos/{id}", UpdateTaskEndpoint).Methods("PUT")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: false,
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":12345", handler ))

}