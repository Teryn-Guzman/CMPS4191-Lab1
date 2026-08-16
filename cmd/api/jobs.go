package main

import (
	"encoding/json"
	"net/http"

	"github.com/Teryn-Guzman/Gatekeeper/internal/data"
	"github.com/Teryn-Guzman/Gatekeeper/internal/validator"
)

func (app *applicationDependencies) createJobHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        ConsumerID string          `json:"consumer_id"`
        JobType    string          `json:"job_type"`
        Payload    json.RawMessage `json:"payload"`
    }
    if err := readJSON(w, r, &input); err != nil {
        badRequestResponse(w, r, err)
        return
    }

    job := &data.Job{
        ConsumerID: input.ConsumerID,
        JobType:    input.JobType,
        Payload:    input.Payload,
    }

    v := validator.New()
    v.Check(job.ConsumerID != "", "consumer_id", "must be provided")
    v.Check(job.JobType != "", "job_type", "must be provided")
    if !v.Valid() {
        failedValidationResponse(w, r, v.Errors)
        return
    }

    if err := app.store.InsertJob(r.Context(), job); err != nil {
        serverErrorResponse(w, r, err)
        return
    }
    _ = writeJSON(w, http.StatusCreated, envelope{"job": job})
}

func (app *applicationDependencies) showJobHandler(w http.ResponseWriter, r *http.Request) {
    id, err := readUUIDParam("id", r)
    if err != nil {
        notFoundResponse(w, r)
        return
    }
    j, err := app.store.GetJobByID(r.Context(), id)
    if err != nil {
        serverErrorResponse(w, r, err)
        return
    }
    _ = writeJSON(w, http.StatusOK, envelope{"job": j})
}

func (app *applicationDependencies) updateJobHandler(w http.ResponseWriter, r *http.Request) {
    id, err := readUUIDParam("id", r)
    if err != nil {
        notFoundResponse(w, r)
        return
    }
    j, err := app.store.GetJobByID(r.Context(), id)
    if err != nil {
        serverErrorResponse(w, r, err)
        return
    }

    var input struct {
        JobType *string         `json:"job_type"`
        Status  *string         `json:"status"`
        Payload *json.RawMessage `json:"payload"`
        Result  *json.RawMessage `json:"result"`
    }
    if err := readJSON(w, r, &input); err != nil {
        badRequestResponse(w, r, err)
        return
    }
    if input.JobType != nil { j.JobType = *input.JobType }
    if input.Status != nil { j.Status = *input.Status }
    if input.Payload != nil { j.Payload = *input.Payload }
    if input.Result != nil { j.Result = *input.Result }

    if err := app.store.UpdateJob(r.Context(), &j); err != nil {
        serverErrorResponse(w, r, err)
        return
    }
    _ = writeJSON(w, http.StatusOK, envelope{"job": j})
}

func (app *applicationDependencies) deleteJobHandler(w http.ResponseWriter, r *http.Request) {
    id, err := readUUIDParam("id", r)
    if err != nil {
        notFoundResponse(w, r)
        return
    }
    if err := app.store.DeleteJob(r.Context(), id); err != nil {
        serverErrorResponse(w, r, err)
        return
    }
    _ = writeJSON(w, http.StatusOK, envelope{"message": "job successfully deleted"})
}
