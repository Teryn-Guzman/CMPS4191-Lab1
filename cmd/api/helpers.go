package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(data)
}

type envelope map[string]interface{}

func readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    if err := dec.Decode(dst); err != nil {
        return err
    }
    // ensure only single JSON value
    if dec.More() {
        return fmt.Errorf("request body must only contain a single JSON object")
    }
    return nil
}

func readUUIDParam(name string, r *http.Request) (string, error) {
    // extract last segment of path
    p := strings.TrimSuffix(r.URL.Path, "/")
    parts := strings.Split(p, "/")
    if len(parts) == 0 {
        return "", fmt.Errorf("missing %s", name)
    }
    id := parts[len(parts)-1]
    if id == "" || id == "consumers" {
        return "", fmt.Errorf("missing %s", name)
    }
    return id, nil
}

// error response helpers
func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
    _ = writeJSON(w, http.StatusInternalServerError, envelope{"error": err.Error()})
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
    _ = writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()})
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
    _ = writeJSON(w, http.StatusNotFound, envelope{"error": "not found"})
}

func failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string][]string) {
    _ = writeJSON(w, http.StatusUnprocessableEntity, envelope{"errors": errors})
}
