package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(cfg mysql.Config) *MySQLStorage {
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MySQL")

	return &MySQLStorage{db: db}
}

func (s *MySQLStorage) Init() (*sql.DB, error) {
	// Initialize the database

	if err := s.CreateProjectsTable(); err != nil {
		return nil, err
	}

	if err := s.CreateUsersTable(); err != nil {
		return nil, err
	}

	if err := s.CreateTasksTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *MySQLStorage) CreateProjectsTable() error {
	// Create the projects table

	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (id)
		)	ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) CreateUsersTable() error {
	// Create the users table

	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			email VARCHAR(255) NOT NULL,
			firstName VARCHAR(255) NOT NULL,
			lastName VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (id),
			UNIQUE KEY (email)
		)	ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	return err
}

func (s *MySQLStorage) CreateTasksTable() error {
	// Create the tasks table

	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			status ENUM('TODO', 'IN_PROGRESS', 'IN_TESTING', 'DONE') NOT NULL DEFAULT 'TODO',
			projectId INT UNSIGNED NOT NULL,
			assignedToID INT UNSIGNED NOT NULL,
			createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

			PRIMARY KEY (id),
            FOREIGN KEY (assignedToID) REFERENCES users(id),
            FOREIGN KEY (projectId) REFERENCES projects(id)
        )	ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)

	if err != nil {
		return err
	}

	return nil
}
