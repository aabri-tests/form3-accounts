package errors

import "fmt"

// ErrBadRequest represents a 400 Bad Request error.
type ErrBadRequest struct {
	Detail string
}

func (e *ErrBadRequest) Error() string {
	return fmt.Sprintf("bad request: %s", e.Detail)
}

// ErrNotFound represents a 404 Not Found error.
type ErrNotFound struct {
	ResourceID string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("resource not found with ID: %s", e.ResourceID)
}

// ErrPermanentFailure represents a permanent failure that should not be retried.
type ErrPermanentFailure struct {
	Detail string
}

func (e *ErrPermanentFailure) Error() string {
	return fmt.Sprintf("permanent failure: %s", e.Detail)
}
