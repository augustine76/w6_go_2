package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // "pending" or "completed"
}

var (
	tasks  = []Task{}
	nextID = 1
	mu     sync.Mutex
)

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	task.ID = nextID
	nextID++
	tasks = append(tasks, task)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllTasks(w, r)
		case http.MethodPost:
			createTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
	//     switch r.Method {
	//     case http.MethodGet:
	//         getTaskByID(w, r)
	//     case http.MethodPut:
	//         updateTask(w, r)
	//     case http.MethodDelete:
	//         deleteTask(w, r)
	//     default:
	//         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//     }
	// })

	http.ListenAndServe(":8080", nil)
}
