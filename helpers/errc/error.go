package errc

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"errors"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func GetValidationErrMgs(err error) interface{} {
	var errV validator.ValidationErrors
	
	if errors.As(err, &errV) {
		return ParseErrMgs(errV)
	}

	log.Println(err.Error())
	return http.StatusText(http.StatusBadRequest)
}

func ParseErrMgs(verr validator.ValidationErrors) []ValidationError {
	errs := []ValidationError{}

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs = append(errs, ValidationError{Field: strings.ToLower(f.Field()), Reason: err})
	}

	return errs
}

func GetGORMErrMgs(err error) *ValidationError {
	var errMg ValidationError

	if pgError, ok := err.(*pgconn.PgError); ok {
		mgs := strings.Replace(pgError.Message, "\""+pgError.ConstraintName+"\"", "", -1)
		errMg.Field = extractString(pgError.ConstraintName)
		errMg.Reason = strings.Trim(mgs, " ")
	}

	return &errMg
}

func extractString(input string) string {
	// Find the last underscore in the input string
	lastIndex := strings.LastIndex(input, "_")

	if lastIndex == -1 {
		// If no underscore is found, return the original string
		return input
	}
	// Extract the substring that follows the last underscore
	extractedString := input[lastIndex+1:]

	return extractedString
}
