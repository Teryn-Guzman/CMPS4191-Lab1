package main

import (
	"net/http"
)

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
    _ = writeJSON(w, http.StatusOK, map[string]interface{}{"status": "available", "version": appVersion})
}
