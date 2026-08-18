package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Teryn-Guzman/Gatekeeper/internal/validator"
)

// createReportHandler generates a consumer activity report for a given time range.
func (app *applicationDependencies) createReportHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConsumerID string    `json:"consumer_id"`
		From       time.Time `json:"from"`
		To         time.Time `json:"to"`
	}

	if err := readJSON(w, r, &input); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.ConsumerID != "", "consumer_id", "must be provided")
	v.Check(!input.From.IsZero(), "from", "must be provided")
	v.Check(!input.To.IsZero(), "to", "must be provided")
	v.Check(input.From.Before(input.To), "from", "must be earlier than to")
	if !v.Valid() {
		failedValidationResponse(w, r, v.Errors)
		return
	}

	app.logger.Info("report generation started", "consumer_id", input.ConsumerID)

	report, err := app.store.GenerateReport(r.Context(), input.ConsumerID, input.From, input.To)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			notFoundResponse(w, r)
		} else {
			serverErrorResponse(w, r, err)
		}
		return
	}

	app.logger.Info("report generation finished", "consumer_id", input.ConsumerID)

	_ = writeJSON(w, http.StatusOK, envelope{"report": report})
}
