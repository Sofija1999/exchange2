package response

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func Ok(w http.ResponseWriter, v interface{}) error {
	return Send(w, http.StatusOK, v)
}

func Created(w http.ResponseWriter, v interface{}) error {
	return Send(w, http.StatusCreated, v)
}

func InternalServerError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func Send(w http.ResponseWriter, code int, v interface{}) error {
	writeHeaders(w)
	w.WriteHeader(code)

	if v == nil {
		return nil
	}

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return errors.Wrap(err, "encode response")
	}

	return nil
}

func writeHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
