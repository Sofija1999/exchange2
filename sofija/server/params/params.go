package params

import (
	"context"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

var (
	ErrNotFound      = errors.New("param not found")
	ErrMissingParams = errors.New("missing params")
)

type key int

var paramsKey key

func From(r *http.Request, key string) (string, bool) {
	params := r.Context().Value(paramsKey).(map[string]string)
	if params == nil {
		return "", false
	}

	param, ok := params[key]
	if !ok {
		return "", false
	}

	return param, true
}

func StringFrom(r *http.Request, key string) (string, error) {
	// sanity checks
	if r.Context().Value(paramsKey) == nil {
		return "", ErrMissingParams
	}
	params := r.Context().Value(paramsKey).(map[string]string)
	if params == nil {
		return "", ErrMissingParams
	}

	param, ok := params[key]
	if !ok {
		return "", ErrNotFound
	}

	return param, nil
}

func IntFrom(r *http.Request, key string) (int, error) {
	param, err := StringFrom(r, key)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(param)
	if err != nil {
		return 0, errors.Wrap(err, "invalid int param")
	}

	return i, nil
}

func WithRequest(r *http.Request, params httprouter.Params) *http.Request {
	ctx := context.WithValue(r.Context(), paramsKey, paramsToMap(params))

	return r.WithContext(ctx)
}

func paramsToMap(params httprouter.Params) map[string]string {
	p := make(map[string]string)
	for _, param := range params {
		p[param.Key] = param.Value
	}

	return p
}
