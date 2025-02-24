package main

// Mocks

type MockStore struct{}

// func (m *MockStore) CreateUser() error {
// 	return nil
// }

func (ms *MockStore) CreateUser(user *User) (*User, error) {
	// Mock implementation
	user.ID = 1
	return user, nil
}

func (m *MockStore) CreateTask(t *Task) (*Task, error) {
	return &Task{}, nil
}

func (m *MockStore) GetTask(id string) (*Task, error) {
	return &Task{}, nil
}
