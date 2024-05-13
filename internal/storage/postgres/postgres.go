package postgres

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type Task struct {
	ID          int
	Title       string
	Description string
}

func New() (*Storage, error) {
	const op = "storage.postgres.New"

	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}

	var (
		dbUser    = os.Getenv("DB_USER")
		dbPass    = os.Getenv("DB_PASS")
		dbName    = os.Getenv("DB_NAME")
		dbSSLMode = os.Getenv("DB_SSLMODE")
		uri       = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUser, dbPass, dbName, dbSSLMode)
	)

	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Create(title, description string) error {
	const op = "storage.postgres.Create"

	_, err := s.db.Exec("INSERT INTO task (title, description) VALUES ($1, $2)", title, description)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Delete(id int) error {
	const op = "storage.postgres.Delete"
	_, err := s.db.Exec("DELETE FROM task WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Update(id int, title, description string) error {
	const op = "storage.postgres.Update"

	_, err := s.db.Exec("UPDATE task set title = $1, description = $2 WHERE id = $3", title, description, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Tasks() ([]Task, error) {
	const op = "storage.postgres.Tasks"

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

func (s *Storage) Task(id int) (Task, error) {
	const op = "storage.postgres.Task"

	var task Task

	row := s.db.QueryRow("SELECT * FROM task WHERE id = $1", id)
	err := row.Scan(&task.ID, &task.Title, &task.Description)
	if err != nil {
		return task, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}
