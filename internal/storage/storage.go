package storage

import "errors"

var (
	ErrIDNotFound    = errors.New("id not found")
	ErrTasksNotFound = errors.New("tasks not found")
	ErrTaskExists    = errors.New("task already exists")
)
