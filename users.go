package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type UserService struct {
	store Store
}

func NewUserService(s Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", s.GetUser).Methods("GET")
}

func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Create a user

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var payload *User

	err = json.Unmarshal(body, &payload)

	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request payload"})
		return
	}

	if err := validateUserPayload(payload); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hashedPassword, err := hashPassword(payload.Password)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to hash password"})
		return
	}

	payload.Password = hashedPassword

	u, err := s.store.CreateUser(payload)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to create user"})
		return
	}

	token, err := createAndSetAuthCookie(u.ID, w)

	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to create token"})
		return
	}

	WriteJSON(w, http.StatusCreated, token)
}

func validateUserPayload(user *User) error {
	if user.Email == "" {
		return nil //errEmailRequired
	}

	if user.FirstName == "" {
		return nil //errFirstNameRequired
	}

	if user.LastName == "" {
		return nil //errLastNameRequired
	}

	if user.Password == "" {
		return nil //errPasswordRequired
	}

	return nil
}

func (s *UserService) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get a user
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	// Create a JWT token

	secret := []byte(Envs.JWTSecret)

	token, err := CreateJWT(secret, id)

	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}
