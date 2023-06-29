package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bloxico/exchange-gateway/sofija/server/response"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

func ReadBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return response.NewValidationError("invalid request body").WithInternal(err)
	}

	validate := validator.New()
	err = validate.Struct(v)
	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			return response.NewValidationError("Error in field: " + strings.ToLower(fieldErr.Field()) +
				". Expected: " + fieldErr.ActualTag() + " " + fieldErr.Param()).WithInternal(err)
		}
	}

	return nil
}

type MissingParamError struct {
	missingParameter string
}

func (e MissingParamError) Error() string {
	return fmt.Sprintf("Missing parameter: %s", e.missingParameter)
}

func QueryParam(r *http.Request, k string, v string) string {
	param := r.URL.Query().Get(k)
	if param == "" {
		return v
	}

	return param
}

func QueryMultipleParam(r *http.Request, k string, v []string) []string {
	params := r.URL.Query().Get(k)
	if params == "" {
		return v
	}
	paramArray := strings.Split(params, ",")

	return paramArray
}

func IntQueryParam(r *http.Request, k string, v int) (int, error) {
	param := QueryParam(r, k, "")

	if param == "" {
		return v, nil
	}

	i, err := strconv.Atoi(param)
	if err != nil {
		return 0, errors.Wrap(err, "parse query param")
	}

	return i, nil
}

func RequiredQueryParam(r *http.Request, k string) (string, error) {
	param := QueryParam(r, k, "")

	if param == "" {
		return "", MissingParamError{k}
	}

	return param, nil
}

func RequiredIntQueryParam(r *http.Request, k string) (int, error) {
	param := QueryParam(r, k, "")

	if param == "" {
		return 0, MissingParamError{k}
	}

	i, err := strconv.Atoi(param)
	if err != nil {
		return 0, errors.Wrap(err, "parse query param")
	}

	return i, nil
}

func BoolQueryParam(r *http.Request, k string, v bool) (bool, error) {
	param := QueryParam(r, k, "")

	if param == "" {
		return v, nil
	}

	b, err := strconv.ParseBool(param)
	if err != nil {
		return false, errors.Wrap(err, "parse query param")
	}

	return b, nil
}
