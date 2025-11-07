package main

import (
	"encoding/json"
	"net/http"
)

// envelope is a lightweight wrapper used to create JSON responses
// with a clear top-level key. Example usage:
//
//	writeJSON(w, http.StatusOK, envelope{"user": u}, nil)
//
// will produce:
//
//	{"user": {...}}
type envelope map[string]interface{}

// writeJSON writes a JSON response with the given status code and optional headers.
//   - Sets "Content-Type: application/json".
//   - Disables HTML escaping to prevent characters like '<' or '>' from being escaped.
//   - Returns an error if encoding fails, so the caller can handle it.
//
// Example:
//
//	data := envelope{"users": list} // or {"user": user}, {"token": token}, etc.
//	if err := writeJSON(w, http.StatusOK, data, nil); err != nil { ... }
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
