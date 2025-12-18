package usecase

import "errors"

var (
	// ErrRepositorySave is returned when there is an error saving to the repository.
	ErrRepositorySave = errors.New("repository save error")
	// ErrDBFindAll is returned when there is an error finding all records in the database.
	ErrDBFindAll  = errors.New("db find all error")
	ErrDBFindByID = errors.New("db find by id error")
	ErrDBDelete   = errors.New("db delete error")
)
