package usecase

import "errors"

var (
	// ErrRepositorySave is returned when there is an error saving to the repository.
	ErrRepositorySave = errors.New("repository save error")
	// ErrDBFindAll is returned when there is an error finding all records in the database.
	ErrDBFindAll  = errors.New("db find all error")
	ErrFindByIDDB = errors.New("FindByID db error")
	ErrDBDelete   = errors.New("db delete error")
)
