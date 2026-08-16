# Gatekeeper - Minimal API for testing migrations

This small server exposes a `/debug` endpoint that returns JSON for `consumers`, `api_keys`, and `jobs` tables.

Prerequisites:
- A Postgres instance with the migrations applied (see `migrations/`)
- `DATABASE_URL` environment variable set (e.g. `postgres://user:pass@localhost:5432/dbname`)

Run locally:

```bash
export DATABASE_URL='postgres://user:pass@localhost:5432/dbname'
cd Gatekeeper
go mod tidy
go run ./cmd/api
```

Then open http://localhost:8080/debug to see JSON output.
