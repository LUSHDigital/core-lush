package lusherr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/LUSHDigital/core-lush/lusherr"
	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/uuid"
)

var inner = fmt.Errorf("this is the inner most error")

func TestInternalError_Error(t *testing.T) {
	test.Equals(t, "internal failure", lusherr.NewInternalError(nil).Error())
	test.Equals(t, "internal failure: inner", lusherr.NewInternalError(fmt.Errorf("inner")).Error())
	test.Equals(t, inner, errors.Unwrap(lusherr.NewInternalError(inner)))
}

func TestUnauthorizedError_Error(t *testing.T) {
	test.Equals(t, "unauthorized", lusherr.NewUnauthorizedError(nil).Error())
	test.Equals(t, "unauthorized: inner", lusherr.NewUnauthorizedError(fmt.Errorf("inner")).Error())
	test.Equals(t, inner, errors.Unwrap(lusherr.NewUnauthorizedError(inner)))
}

func TestValidationError_Error(t *testing.T) {
	test.Equals(t, "validation failed for \"name\" on \"user\"", lusherr.NewValidationError("user", "name", nil).Error())
	test.Equals(t, "validation failed for \"name\" on \"user\": inner", lusherr.NewValidationError("user", "name", fmt.Errorf("inner")).Error())
	test.Equals(t, inner, errors.Unwrap(lusherr.NewValidationError("user", "name", inner)))
}

func TestDatabaseQueryError_Error(t *testing.T) {
	test.Equals(t, "database query failed", lusherr.NewDatabaseQueryError("SELECT * FROM user", nil).Error())
	test.Equals(t, "database query failed: inner", lusherr.NewDatabaseQueryError("SELECT * FROM user", fmt.Errorf("inner")).Error())
	test.Equals(t, inner, errors.Unwrap(lusherr.NewDatabaseQueryError("", inner)))
}

func TestNotFoundError_Error(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	test.Equals(t, fmt.Sprintf("cannot find \"user\" (%s)", id), lusherr.NewNotFoundError("user", id, nil).Error())
	test.Equals(t, fmt.Sprintf("cannot find \"user\" (%s): inner", id), lusherr.NewNotFoundError("user", id, fmt.Errorf("inner")).Error())
	test.Equals(t, inner, errors.Unwrap(lusherr.NewNotFoundError("", "", inner)))
}

func TestNotAllowedError_Error(t *testing.T) {
	test.Equals(t, "action not allowed", lusherr.NewNotAllowedError(nil).Error())
	test.Equals(t, "action not allowed: inner", lusherr.NewNotAllowedError(fmt.Errorf("inner")).Error())
	test.Equals(t, inner, errors.Unwrap(lusherr.NewNotAllowedError(inner)))
}
