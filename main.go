package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Assignment struct {
	ID             string    `json:"id"`
	CourseName     string    `json:"course_name"`
	CourseCode     string    `json:"course_code"`
	AssignmentWeek int       `json:"assignment_week"`
	AssignmentDue  time.Time `json:"assignment_due"`
}

var assignments = make(map[string]Assignment)

func calculateHoursLeft(due time.Time) float64 {
	duration := time.Until(due)
	return duration.Hours()
}

func createAssignment(w http.ResponseWriter, r *http.Request) {
	var assignment Assignment
	if err := json.NewDecoder(r.Body).Decode(&assignment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	assignment.ID = uuid.New().String()
	assignments[assignment.ID] = assignment
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assignment)
}

func getAssignments(w http.ResponseWriter, r *http.Request) {
	var list []map[string]interface{}
	for _, assignment := range assignments {
		assignmentData := map[string]interface{}{
			"id":              assignment.ID,
			"course_name":     assignment.CourseName,
			"course_code":     assignment.CourseCode,
			"assignment_week": assignment.AssignmentWeek,
			"assignment_due":  assignment.AssignmentDue,
			"hours_left":      calculateHoursLeft(assignment.AssignmentDue),
		}
		list = append(list, assignmentData)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func getAssignment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	assignment, exists := assignments[params["id"]]
	if !exists {
		http.Error(w, "Assignment not found", http.StatusNotFound)
		return
	}

	assignmentData := map[string]interface{}{
		"id":              assignment.ID,
		"course_name":     assignment.CourseName,
		"course_code":     assignment.CourseCode,
		"assignment_week": assignment.AssignmentWeek,
		"assignment_due":  assignment.AssignmentDue,
		"hours_left":      calculateHoursLeft(assignment.AssignmentDue),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assignmentData)
}

func updateAssignment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, exists := assignments[params["id"]]
	if !exists {
		http.Error(w, "Assignment not found", http.StatusNotFound)
		return
	}

	var updated Assignment
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updated.ID = params["id"]
	assignments[params["id"]] = updated
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func deleteAssignment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, exists := assignments[params["id"]]
	if !exists {
		http.Error(w, "Assignment not found", http.StatusNotFound)
		return
	}
	delete(assignments, params["id"])
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/assignments", createAssignment).Methods("POST")
	router.HandleFunc("/assignments", getAssignments).Methods("GET")
	router.HandleFunc("/assignments/{id}", getAssignment).Methods("GET")
	router.HandleFunc("/assignments/{id}", updateAssignment).Methods("PUT")
	router.HandleFunc("/assignments/{id}", deleteAssignment).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}
