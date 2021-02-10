package rest_errors

import (
	"errors"
	"net/http"
)

// RestErr is a rest error struct
type RestErr struct {
	Message string        `json:"message"`
	Status  int           `json:"status"`
	Error   string        `json:"error"`
	Causes  []interface{} `json:"causes"` // causes of the error
}

// NewError returns a general message of the error.
// NewError is largely used to send a vague description back to an external caller.
func NewError(msg string) error {
	return errors.New(msg)
}

// NewBadRequestError returns a status bad request
func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

// NewNotFoundError returns a status not found
func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

// NewUnauthorizedError returns a rest error for unauthorized user
func NewUnauthorizedError() *RestErr {
	return &RestErr{
		Message: "unauthorized user",
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized_user_error",
	}
}

// NewInternalServerError returns an internal server error
func NewInternalServerError(message string, err error) *RestErr {
	result := &RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
	if err != nil {
		result.Causes = append(result.Causes, err.Error())
	}
	return result
}
