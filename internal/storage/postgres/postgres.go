package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func New() (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", "user=zeze password=123qwe123 dbname=postgres sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateTask(id int, title, description string) error {
	const op = "storage.postgres.CreateTask"

	_, err := s.db.Exec("INSERT INTO task (id, title, description) VALUES ($1, $2, $3)", id, title, description)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteTask(id int) error {
	const op = "storage.postgres.DeleteTask"

	_, err := s.db.Exec("DELETE FROM task WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateTask(id int, title, description string) error {
	const op = "storage.postgres.UpdateTask"

	_, err := s.db.Exec("UPDATE task set title = $1, description = $2 WHERE id = $3", title, description, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetTasks() ([]Task, error) {
	const op = "storage.postgres.GetTasks"

	rows, err := s.db.Query("SELECT * FROM task")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		if err = rows.Scan(&task.ID, &task.Title, &task.Description); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) GetTask(id int) (Task, error) {
	const op = "storage.postgres.GetTask"

	var task Task

	row := s.db.QueryRow("SELECT * FROM task WHERE id = $1", id)
	err := row.Scan(&task.ID, &task.Title, &task.Description)
	if err != nil {
		return task, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}
