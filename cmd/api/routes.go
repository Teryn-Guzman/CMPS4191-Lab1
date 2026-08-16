package main

import (
	"net/http"
)

func routes(app *applicationDependencies) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        consumers, err := app.store.ListConsumers(ctx, 100)
        if err != nil {
            serverErrorResponse(w, r, err)
            return
        }
        apiKeys, err := app.store.ListAPIKeys(ctx, 100)
        if err != nil {
            serverErrorResponse(w, r, err)
            return
        }
        jobs, err := app.store.ListJobs(ctx, 100)
        if err != nil {
            serverErrorResponse(w, r, err)
            return
        }

        resp := map[string]interface{}{
            "consumers": consumers,
            "api_keys":  apiKeys,
            "jobs":      jobs,
        }
        _ = writeJSON(w, http.StatusOK, resp)
    })

    // consumers CRUD
    mux.HandleFunc("/consumers", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            // list
            consumers, err := app.store.ListConsumers(r.Context(), 100)
            if err != nil {
                serverErrorResponse(w, r, err)
                return
            }
            _ = writeJSON(w, http.StatusOK, envelope{"consumers": consumers})
        case http.MethodPost:
            app.createConsumerHandler(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/consumers/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            app.showConsumerHandler(w, r)
        case http.MethodPatch, http.MethodPut:
            app.updateConsumerHandler(w, r)
        case http.MethodDelete:
            app.deleteConsumerHandler(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    // existing endpoints
    mux.HandleFunc("/api_keys", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            keys, err := app.store.ListAPIKeys(r.Context(), 100)
            if err != nil {
                serverErrorResponse(w, r, err)
                return
            }
            _ = writeJSON(w, http.StatusOK, envelope{"api_keys": keys})
        case http.MethodPost:
            app.createAPIKeyHandler(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/api_keys/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            app.showAPIKeyHandler(w, r)
        case http.MethodPatch, http.MethodPut:
            app.updateAPIKeyHandler(w, r)
        case http.MethodDelete:
            app.deleteAPIKeyHandler(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            jobs, err := app.store.ListJobs(r.Context(), 100)
            if err != nil {
                serverErrorResponse(w, r, err)
                return
            }
            _ = writeJSON(w, http.StatusOK, envelope{"jobs": jobs})
        case http.MethodPost:
            app.createJobHandler(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            app.showJobHandler(w, r)
        case http.MethodPatch, http.MethodPut:
            app.updateJobHandler(w, r)
        case http.MethodDelete:
            app.deleteJobHandler(w, r)
        default:
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/healthz", healthcheckHandler)

    return LoggerMiddleware(app.logger)(mux)
}
