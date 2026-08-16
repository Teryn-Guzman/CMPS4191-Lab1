package main

import (
	"net/http"

	"github.com/Teryn-Guzman/Gatekeeper/internal/data"
	"github.com/Teryn-Guzman/Gatekeeper/internal/validator"
)

// createConsumerHandler reads JSON input to create a consumer record, then returns JSON of the created record.
func (app *applicationDependencies) createConsumerHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }

    if err := readJSON(w, r, &input); err != nil {
        badRequestResponse(w, r, err)
        return
    }

    consumer := &data.Consumer{
        Name:  input.Name,
        Email: input.Email,
    }

    v := validator.New()
    v.Check(consumer.Name != "", "name", "must be provided")
    validator.ValidateContactEmail(v, consumer.Email)
    if !v.Valid() {
        failedValidationResponse(w, r, v.Errors)
        return
    }

    if err := app.store.InsertConsumer(r.Context(), consumer); err != nil {
        serverErrorResponse(w, r, err)
        return
    }

    _ = writeJSON(w, http.StatusCreated, envelope{"consumer": consumer})
}

// showConsumerHandler reads the UUID in the URL path, then returns JSON of the matching consumer record.
func (app *applicationDependencies) showConsumerHandler(w http.ResponseWriter, r *http.Request) {
    id, err := readUUIDParam("id", r)
    if err != nil {
        notFoundResponse(w, r)
        return
    }

    consumer, err := app.store.GetConsumerByID(r.Context(), id)
    if err != nil {
        serverErrorResponse(w, r, err)
        return
    }

    _ = writeJSON(w, http.StatusOK, envelope{"consumer": consumer})
}

// updateConsumerHandler updates the consumer record matching UUID in the URL path, then returns JSON of the updated record.
func (app *applicationDependencies) updateConsumerHandler(w http.ResponseWriter, r *http.Request) {
    id, err := readUUIDParam("id", r)
    if err != nil {
        notFoundResponse(w, r)
        return
    }

    consumer, err := app.store.GetConsumerByID(r.Context(), id)
    if err != nil {
        serverErrorResponse(w, r, err)
        return
    }

    var input struct {
        Name   *string `json:"name"`
        Email  *string `json:"email"`
        Status *string `json:"status"`
    }

    if err := readJSON(w, r, &input); err != nil {
        badRequestResponse(w, r, err)
        return
    }

    if input.Name != nil {
        consumer.Name = *input.Name
    }
    if input.Email != nil {
        consumer.Email = *input.Email
    }
    if input.Status != nil {
        consumer.Status = *input.Status
    }

    v := validator.New()
    if !validator.ValidateConsumer(v, consumer) {
        failedValidationResponse(w, r, v.Errors)
        return
    }

    if err := app.store.UpdateConsumer(r.Context(), &consumer); err != nil {
        serverErrorResponse(w, r, err)
        return
    }

    _ = writeJSON(w, http.StatusOK, envelope{"consumer": consumer})
}

// deleteConsumerHandler deletes the consumer record matching UUID in the URL path.
func (app *applicationDependencies) deleteConsumerHandler(w http.ResponseWriter, r *http.Request) {
    id, err := readUUIDParam("id", r)
    if err != nil {
        notFoundResponse(w, r)
        return
    }

    if err := app.store.DeleteConsumer(r.Context(), id); err != nil {
        serverErrorResponse(w, r, err)
        return
    }

    _ = writeJSON(w, http.StatusOK, envelope{"message": "consumer successfully deleted"})
}
