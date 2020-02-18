package lusherr

import (
	"fmt"
	"runtime"
)

// NewInternalError builds a generic error.
// e.g. Trying to generate a random UUID, but the generation failed.
func NewInternalError(inner error) error {
	return InternalError{
		frame: frame(1),
		inner: inner,
	}
}

// InternalError can be used to wrap any error.
// e.g. Trying to generate a random UUID, but the generation failed.
type InternalError struct {
	frame runtime.Frame
	inner error
}

func (e InternalError) Error() string {
	if e.inner == nil {
		return fmt.Sprintf("internal failure")
	}
	return fmt.Sprintf("internal failure: %v", e.inner)
}

// Unwrap the inner error.
func (e InternalError) Unwrap() error {
	return e.inner
}

// Locate the frame of the error.
func (e InternalError) Locate() runtime.Frame {
	return e.frame
}

// Pin the error to a caller frame.
func (e InternalError) Pin(frame runtime.Frame) error {
	e.frame = frame
	return e
}

// NewUnauthorizedError builds a new unauthorized error.
// e.g. Someone tried to access something they were not allowed to according to a permission policy.
func NewUnauthorizedError(inner error) error {
	return UnauthorizedError{
		frame: frame(1),
		inner: inner,
	}
}

// UnauthorizedError should be used when an action is performed by a user that they don't have permission to do.
// e.g. Someone tried to access something they were not allowed to according to a permission policy.
type UnauthorizedError struct {
	frame runtime.Frame
	inner error
}

func (e UnauthorizedError) Error() string {
	if e.inner == nil {
		return fmt.Sprintf("unauthorized")
	}
	return fmt.Sprintf("unauthorized: %v", e.inner)
}

// Unwrap the inner error.
func (e UnauthorizedError) Unwrap() error {
	return e.inner
}

// Locate the frame of the error.
func (e UnauthorizedError) Locate() runtime.Frame {
	return e.frame
}

// Pin the error to a caller frame.
func (e UnauthorizedError) Pin(frame runtime.Frame) error {
	e.frame = frame
	return e
}

// NewValidationError builds an error for failing to validate a field.
// e.g. Someone set the name field for a user to be empty, but the validation requires it to be present.
func NewValidationError(entity, field string, inner error) error {
	return ValidationError{
		frame:  frame(1),
		inner:  inner,
		Entity: entity,
		Field:  field,
	}
}

// ValidationError should be used to detail what user generated information is incorrect and why.
// e.g. Someone set the name field for a user to be empty, but the validation requires it to be present.
type ValidationError struct {
	Entity, Field string
	frame         runtime.Frame
	inner         error
}

func (e ValidationError) Error() string {
	if e.inner == nil {
		return fmt.Sprintf("validation failed for %q on %q", e.Field, e.Entity)
	}
	return fmt.Sprintf("validation failed for %q on %q: %v", e.Field, e.Entity, e.inner)
}

// Unwrap the inner error.
func (e ValidationError) Unwrap() error {
	return e.inner
}

// Locate the frame of the error.
func (e ValidationError) Locate() runtime.Frame {
	return e.frame
}

// Pin the error to a caller frame.
func (e ValidationError) Pin(frame runtime.Frame) error {
	e.frame = frame
	return e
}

// NewDatabaseQueryError builds an error for a failed database query.
// e.g. Trying to query the database, but the database rejects the query.
func NewDatabaseQueryError(query string, inner error) error {
	return DatabaseQueryError{
		frame: frame(1),
		inner: inner,
		Query: query,
	}
}

// DatabaseQueryError should be used to provide detail about a failed database query.
// e.g. Trying to query the database, but the database rejects the query.
type DatabaseQueryError struct {
	Query string
	frame runtime.Frame
	inner error
}

func (e DatabaseQueryError) Error() string {
	if e.inner == nil {
		return fmt.Sprintf("database query failed")
	}
	return fmt.Sprintf("database query failed: %v", e.inner)
}

// Unwrap the inner error.
func (e DatabaseQueryError) Unwrap() error {
	return e.inner
}

// Locate the frame of the error.
func (e DatabaseQueryError) Locate() runtime.Frame {
	return e.frame
}

// Pin the error to a caller frame.
func (e DatabaseQueryError) Pin(frame runtime.Frame) error {
	e.frame = frame
	return e
}

// NewNotFoundError builds an error for an entity that cannot be found.
// e.g. Someone tries to retrieve a user, but the user for the given ID does not exist in the database.
func NewNotFoundError(entity string, identifier interface{}, inner error) error {
	return NotFoundError{
		frame:      frame(1),
		inner:      inner,
		Entity:     entity,
		Identifier: identifier,
	}
}

// NotFoundError should be used when an entity cannot be found.
// e.g. Someone tries to retrieve a user, but the user for the given ID does not exist in the database.
type NotFoundError struct {
	Entity     string
	Identifier interface{}
	frame      runtime.Frame
	inner      error
}

func (e NotFoundError) Error() string {
	if e.inner == nil {
		return fmt.Sprintf("cannot find %q (%v)", e.Entity, e.Identifier)
	}
	return fmt.Sprintf("cannot find %q (%v): %v", e.Entity, e.Identifier, e.inner)
}

// Unwrap the inner error.
func (e NotFoundError) Unwrap() error {
	return e.inner
}

// Locate the frame of the error.
func (e NotFoundError) Locate() runtime.Frame {
	return e.frame
}

// Pin the error to a caller frame.
func (e NotFoundError) Pin(frame runtime.Frame) error {
	e.frame = frame
	return e
}

// NewNotAllowedError builds an error for when a certain action is not allowed.
// e.g. Someone tries to delete something, but the record has been marked as permanenet.
func NewNotAllowedError(inner error) error {
	return NotAllowedError{
		frame: frame(1),
		inner: inner,
	}
}

// NotAllowedError should be used when an certain action is not allowed.
// e.g. Someone tries to delete something, but the record has been marked as permanent.
type NotAllowedError struct {
	frame runtime.Frame
	inner error
}

func (e NotAllowedError) Error() string {
	if e.inner == nil {
		return fmt.Sprintf("action not allowed")
	}
	return fmt.Sprintf("action not allowed: %v", e.inner)
}

// Unwrap the inner error.
func (e NotAllowedError) Unwrap() error {
	return e.inner
}

// Locate the frame of the error.
func (e NotAllowedError) Locate() runtime.Frame {
	return e.frame
}

// Pin the error to a caller frame.
func (e NotAllowedError) Pin(frame runtime.Frame) error {
	e.frame = frame
	return e
}
