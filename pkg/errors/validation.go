package errors

import (
	"reflect"

	"github.com/go-playground/validator"
)

type ValidationError interface {
	ResponseError
	ValidationErrors() []ValidationErrorField
}

type ValidationErrorField struct {
	Field         string `json:"field,omitempty"`
	Type          string `json:"type,omitempty"`
	Tag           string `json:"tag,omitempty"`
	CriticalValue string `json:"critical_value,omitempty"`
}

type ApiValidationError struct {
	APIResponseError
	ApiValidationErrors []ValidationErrorField
}

type ValidationProblemJson struct {
	ProblemJson
	ValidationErrors []ValidationErrorField `json:"validation_errors,omitempty"`
}

func (v *ApiValidationError) ValidationErrors() []ValidationErrorField {
	return v.ApiValidationErrors
}

func convertReflectTypeForAPI(kind reflect.Kind) string {
	switch kind {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Array, reflect.Slice:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	case reflect.String:
		return "string"
	default:
		return "unknown"
	}
}

func MakeDataValidationErrAPIResponse(title, detail string, code, status int) func(errs ...error) ValidationError {
	return func(errs ...error) ValidationError {
		var err error
		var validationErrors []ValidationErrorField

		r := &ApiValidationError{
			APIResponseError: APIResponseError{
				Detail: detail,
				Title:  title,
				Status: status,
				Code:   code,
			},
		}

		if len(errs) > 0 {
			err = errs[0]
			r.Err = err
		}

		if fieldErrors, ok := err.(validator.ValidationErrors); ok {
			validationErrors = make([]ValidationErrorField, 0, len(fieldErrors))

			for _, fieldErr := range fieldErrors {
				validationErrors = append(validationErrors, ValidationErrorField{
					Field:         fieldErr.Field(),
					Type:          convertReflectTypeForAPI(fieldErr.Kind()),
					Tag:           fieldErr.Tag(),
					CriticalValue: fieldErr.Param(),
				})
			}

			r.ApiValidationErrors = validationErrors
		}

		return r
	}
}
