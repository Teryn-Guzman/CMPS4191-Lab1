# CMPS4191 • LABORATORY 1

## What is included

- CRUD endpoints for consumers
- CRUD endpoints for API keys
- CRUD endpoints for jobs
- `/healthz` health check
- `/debug` endpoint that returns the current data in JSON
- `/reports` endpoint to generate a consumer activity summary

## Prerequisites

- Go 1.22+
- PostgreSQL running locally
- `psql` and `migrate` available in your shell
- A database named `adv_web` with the credentials in `.envrc`

The project expects this DSN by default:

```bash
postgres://adv_web:password@localhost:5432/adv_web?sslmode=disable
```

## 1) Clone and enter the project

```bash
cd /path/to/Gatekeeper
```

## 2) Load environment variables

The repository includes a `.envrc` file that sets the app port and database DSN.

```bash
source .envrc
```

You should now have:

```bash
echo $PORT
echo $ENVIRONMENT
echo $GATEKEEPER_DB_DSN
```

## 3) Make sure PostgreSQL is running

If your local Postgres server is not up, start it first. Then confirm the database exists.

```bash
psql "postgres://adv_web:password@localhost:5432/postgres" -c "SELECT 1;"
```

If needed, create the database:

```bash
createdb -h localhost -U postgres adv_web
```

## 4) Apply the database migrations

```bash
make db/migrations/up
```

This reads the SQL files under `migrations/` and applies them to the configured DSN.

## 5) Run the API

```bash
make run/api
```

The server starts on port `4000` by default.

## 6) Validate the API is running

Health check:

```bash
curl -i http://localhost:4000/healthz
```

Expected result: HTTP `200 OK`.

## Example endpoints

List consumers:

```bash
curl -sS http://localhost:4000/consumers | jq
```

Create a consumer:

```bash
curl -sS -X POST http://localhost:4000/consumers \
  -H "Content-Type: application/json" \
  -d '{"name":"Acme Inc","email":"alice@example.com"}' | jq
```

Generate a report:

```bash
curl -sS -X POST http://localhost:4000/reports \
  -H "Content-Type: application/json" \
  -d '{
    "consumer_id":"0198f000-0000-7000-8000-000000000001",
    "from":"2025-01-01T00:00:00Z",
    "to":"2026-12-31T23:59:59Z"
  }' | jq
```

Debug endpoint:

```bash
curl -sS http://localhost:4000/debug | jq
```

## Useful commands

Install Go dependencies:

```bash
go mod tidy
```

Run directly without Make:

```bash
# prefer using the env var set by .envrc
go run ./cmd/api -port=4000 -env=development -db-dsn="$GATEKEEPER_DB_DSN"
```

Connect to the database:

```bash
make db/psql
```