package main

import (
	"net/http"

	"github.com/Teryn-Guzman/Gatekeeper/internal/data"
	"github.com/Teryn-Guzman/Gatekeeper/internal/validator"
)

func (app *applicationDependencies) createAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        ConsumerID string  `json:"consumer_id"`
        KeyPrefix  string  `json:"key_prefix"`
        Status     *string `json:"status"`
    }
    if err := readJSON(w, r, &input); err != nil {
        badRequestResponse(w, r, err)
        return
    }

    key := &data.APIKey{
        ConsumerID: input.ConsumerID,
        KeyPrefix:  input.KeyPrefix,
    }
    if input.Status != nil {
        key.Status = *input.Status
    }

    v := validator.New()
    v.Check(key.ConsumerID != "", "consumer_id", "must be provided")
    if !v.Valid() {
        failedValidationResponse(w, r, v.Errors)
        return
    }

    if err := app.store.InsertAPIKey(r.Context(), key); err != nil {
        serverErrorResponse(w, r, err)
        return
    }
    _ = writeJSON(w, http.StatusCreated, envelope{"api_key": key})
}

func (app *applicationDependencies) showAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
    id, err := readUUIDParam("id", r)
    if err != nil {
        notFoundResponse(w, r)
        return
    }
    k, err := app.store.GetAPIKeyByID(r.Context(), id)
    if err != nil {
        serverErrorResponse(w, r, err)
        return
    }
    _ = writeJSON(w, http.StatusOK, envelope{"api_key": k})
}

func (app *applicationDependencies) updateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
    id, err := readUUIDParam("id", r)
    if err != nil {
        notFoundResponse(w, r)
        return
    }
    k, err := app.store.GetAPIKeyByID(r.Context(), id)
    if err != nil {
        serverErrorResponse(w, r, err)
        return
    }

    var input struct {
        KeyPrefix *string `json:"key_prefix"`
        Status    *string `json:"status"`
    }
    if err := readJSON(w, r, &input); err != nil {
        badRequestResponse(w, r, err)
        return
    }
    if input.KeyPrefix != nil { k.KeyPrefix = *input.KeyPrefix }
    if input.Status != nil { k.Status = *input.Status }

    v := validator.New()
    if !v.Valid() {
        failedValidationResponse(w, r, v.Errors)
        return
    }

    if err := app.store.UpdateAPIKey(r.Context(), &k); err != nil {
        serverErrorResponse(w, r, err)
        return
    }
    _ = writeJSON(w, http.StatusOK, envelope{"api_key": k})
}

func (app *applicationDependencies) deleteAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
    id, err := readUUIDParam("id", r)
    if err != nil {
        notFoundResponse(w, r)
        return
    }
    if err := app.store.DeleteAPIKey(r.Context(), id); err != nil {
        serverErrorResponse(w, r, err)
        return
    }
    _ = writeJSON(w, http.StatusOK, envelope{"message": "api_key successfully deleted"})
}
