package errors

import (
	stderrors "errors"
	"net/http"
)

type ResponseError interface {
	ErrorDetail() string
	ErrorTitle() string
	ErrorStatus() int
	Error() string
	Cause() error
}

type ProblemJson struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}
type Frame uintptr
type StackTrace []Frame

const (
	ErrCodeBaseInternalServer       = 1000
	ErrCodeBaseBadRequest           = 1001
	ErrCodeBaseBadRequestValidation = 1002
	ErrCodeBaseResourceNotFound     = 1003
	ErrUnsupportedMediaType         = 1004
	ErrNotFound                     = 1005
)

var (
	ErrInternalServer = MakeErrAPIResponse(
		"ErrInternalServer",
		"Unknown error",
		ErrCodeBaseInternalServer,
		http.StatusInternalServerError,
	)

	ErrBadRequest = MakeErrAPIResponse(
		"ErrBadRequest",
		"Unsupported request",
		ErrCodeBaseBadRequest,
		http.StatusBadRequest,
	)

	ErrValidation = func(err string) ResponseError {
		return MakeErrAPIResponse(
			"ErrBadRequestValidation",
			"Validation error: "+err,
			ErrCodeBaseBadRequestValidation,
			http.StatusBadRequest,
		)()
	}

	ErrDataValidation = MakeDataValidationErrAPIResponse(
		"ErrBadRequestValidation",
		"Request validation has failed",
		ErrCodeBaseBadRequestValidation,
		http.StatusBadRequest,
	)

	ErrRouteNotFound = MakeErrAPIResponse(
		"ErrRouteNotFound",
		"Route not found",
		ErrCodeBaseResourceNotFound,
		http.StatusNotFound,
	)

	ErrMethodNotAllowed = MakeErrAPIResponse(
		"ErrMethodNotAllowed",
		"The method you requested is not allowed",
		0,
		http.StatusMethodNotAllowed,
	)

	ErrUnsupportedContentType = MakeErrAPIResponse(
		"ErrUnsupportedMediaType",
		"Unable to process this content type",
		ErrUnsupportedMediaType,
		http.StatusUnsupportedMediaType,
	)

	ErrResourceNotFound = MakeErrAPIResponse(
		"ErrResourceNotFound",
		"The requested resource was not found",
		ErrNotFound,
		http.StatusNotFound,
	)
)

type APIResponseError struct {
	Detail, Title string
	Code, Status  int
	Err           error
}

func (e *APIResponseError) Error() string {
	errStr := e.String()
	if e.Err != nil {
		errStr = errStr + ". " + e.Err.Error()
	}
	return errStr
}

func (e *APIResponseError) Cause() error {
	if e.Err != nil {
		return e.Err
	}
	return nil
}

func (e *APIResponseError) Unwrap() error {
	return e.Cause()
}

func (e *APIResponseError) String() string {
	return e.Title + " (" + e.Detail + ")"
}

func (e *APIResponseError) ErrorDetail() string {
	return e.Detail
}

func (e *APIResponseError) ErrorTitle() string {
	return e.Title
}

func (e *APIResponseError) ErrorStatus() int {
	return e.Status
}

// Tries to find deepest error in errors chain that contains stack
func DeepestErrorWithStack(err error) (error, StackTrace) {
	if err == nil {
		return nil, nil
	}

	// extract all wrapped errors to slice
	errors := []error{err}
	for cause := stderrors.Unwrap(err); cause != nil; cause = stderrors.Unwrap(cause) {
		errors = append(errors, cause)
	}

	// find innermost error with stack
	type stackTracer interface {
		StackTrace() StackTrace
	}

	for i := len(errors) - 1; i >= 0; i-- {
		if errWithStack, ok := errors[i].(stackTracer); ok {
			return errors[i], errWithStack.StackTrace()
		}
	}

	return nil, nil
}

func MakeErrAPIResponse(title, detail string, code, status int) func(err ...error) ResponseError {
	return func(err ...error) ResponseError {
		r := &APIResponseError{
			Detail: detail,
			Title:  title,
			Status: status,
			Code:   code,
		}

		if len(err) > 0 {
			r.Err = err[0]
		}

		return r
	}
}
