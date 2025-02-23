package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")

type TasksService struct {
	store Store
}

func NewTasksService(s Store) *TasksService {
	return &TasksService{store: s}
}

func (s *TasksService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", s.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", s.GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", s.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", s.DeleteTask).Methods("DELETE")
}

func (s *TasksService) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Create a task

	body, err := io.ReadAll(r.Body)

	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload!"})
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var task *Task

	err = json.Unmarshal(body, &task)

	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request payload!"})
		return
	}

	if err := validateTaskPayload(task); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	t, err := s.store.CreateTask(task)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating task"})
		return
	}

	WriteJSON(w, http.StatusCreated, t)
}

func (s *TasksService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	if id == "" {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "ID is required"})
		return
	}

	t, err := s.store.GetTask(id)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Task not found"})
		return
	}

	WriteJSON(w, http.StatusOK, t)
}

func validateTaskPayload(task *Task) error {
	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectID == 0 {
		return errProjectIDRequired
	}

	if task.AssignedToID == 0 {
		return errUserIDRequired
	}

	return nil
}

func (s *TasksService) GetTask(w http.ResponseWriter, r *http.Request) {
	// Get a task
}

func (s *TasksService) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Update a task
}

func (s *TasksService) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Delete a task
}
