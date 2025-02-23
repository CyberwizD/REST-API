package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateTask(t *testing.T) {
	ms := &MockStore{}

	service := NewTasksService(ms)

	t.Run("Return an error if 'name' is empty", func(t *testing.T) {
		payload := &Task{
			Name: "",
		}

		b, err := json.Marshal(payload)

		if err != nil {
			t.Fatalf("Error while marshalling payload: %v", err)
		}

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(b))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router := mux.NewRouter()

		router.HandleFunc("/tasks", service.CreateTask).Methods("POST")

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Error("Invalid Status Code, Fail Error")
		}
	})

	t.Run("Create a task", func(t *testing.T) {
		payload := &Task{
			Name:         "Creating a REST API in GoLang",
			ProjectID:    1,
			AssignedToID: 42,
		}

		b, err := json.Marshal(payload)

		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router := mux.NewRouter()

		router.HandleFunc("/tasks", service.CreateTask)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetTask(t *testing.T) {
	ms := &MockStore{}

	service := NewTasksService(ms)

	t.Run("Get a task", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/tasks/42", nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router := mux.NewRouter()

		router.HandleFunc("/tasks/{id}", service.handleGetTask)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Invalid Status Code: Error")
		}
	})
}
