package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/Teryn-Guzman/Gatekeeper/internal/data"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct{
    db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Store {
    return &Store{db: db}
}

func (s *Store) ListConsumers(ctx context.Context, limit int) ([]data.Consumer, error) {
    rows, err := s.db.Query(ctx, `SELECT id, name, email, status, version, created_at, updated_at FROM consumers ORDER BY created_at DESC LIMIT $1`, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var out []data.Consumer
    for rows.Next() {
        var c data.Consumer
        if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Status, &c.Version, &c.CreatedAt, &c.UpdatedAt); err != nil {
            return nil, err
        }
        out = append(out, c)
    }
    return out, nil
}

func (s *Store) GetConsumerByID(ctx context.Context, id string) (data.Consumer, error) {
    var c data.Consumer
    row := s.db.QueryRow(ctx, `SELECT id, name, email, status, version, created_at, updated_at FROM consumers WHERE id = $1`, id)
    err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Status, &c.Version, &c.CreatedAt, &c.UpdatedAt)
    if err != nil {
        return data.Consumer{}, err
    }
    return c, nil
}

func (s *Store) InsertConsumer(ctx context.Context, c *data.Consumer) error {
    row := s.db.QueryRow(ctx, `INSERT INTO consumers (name, email) VALUES ($1, $2) RETURNING id, version, created_at, updated_at`, c.Name, c.Email)
    return row.Scan(&c.ID, &c.Version, &c.CreatedAt, &c.UpdatedAt)
}

func (s *Store) UpdateConsumer(ctx context.Context, c *data.Consumer) error {
    _, err := s.db.Exec(ctx, `UPDATE consumers SET name=$1, email=$2, status=$3, version = version + 1, updated_at = now() WHERE id=$4`, c.Name, c.Email, c.Status, c.ID)
    return err
}

func (s *Store) DeleteConsumer(ctx context.Context, id string) error {
    _, err := s.db.Exec(ctx, `DELETE FROM consumers WHERE id=$1`, id)
    return err
}

func (s *Store) ListAPIKeys(ctx context.Context, limit int) ([]data.APIKey, error) {
    rows, err := s.db.Query(ctx, `SELECT id, consumer_id, key_prefix, status, last_used_at, expires_at, created_at FROM api_keys ORDER BY created_at DESC LIMIT $1`, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []data.APIKey
    for rows.Next() {
        var a data.APIKey
        if err := rows.Scan(&a.ID, &a.ConsumerID, &a.KeyPrefix, &a.Status, &a.LastUsedAt, &a.ExpiresAt, &a.CreatedAt); err != nil {
            return nil, err
        }
        out = append(out, a)
    }
    return out, nil
}

func (s *Store) ListJobs(ctx context.Context, limit int) ([]data.Job, error) {
    rows, err := s.db.Query(ctx, `SELECT id, consumer_id, job_type, status, payload, result, error_message, started_at, completed_at, created_at FROM jobs ORDER BY created_at DESC LIMIT $1`, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var out []data.Job
    for rows.Next() {
        var j data.Job
        var payload []byte
        var result []byte
        if err := rows.Scan(&j.ID, &j.ConsumerID, &j.JobType, &j.Status, &payload, &result, &j.ErrorMessage, &j.StartedAt, &j.CompletedAt, &j.CreatedAt); err != nil {
            return nil, err
        }
        j.Payload = json.RawMessage(payload)
        j.Result = json.RawMessage(result)
        out = append(out, j)
    }
    return out, nil
}

// APIKeys CRUD
func (s *Store) GetAPIKeyByID(ctx context.Context, id string) (data.APIKey, error) {
    var a data.APIKey
    row := s.db.QueryRow(ctx, `SELECT id, consumer_id, key_prefix, status, last_used_at, expires_at, created_at FROM api_keys WHERE id=$1`, id)
    err := row.Scan(&a.ID, &a.ConsumerID, &a.KeyPrefix, &a.Status, &a.LastUsedAt, &a.ExpiresAt, &a.CreatedAt)
    if err != nil {
        return data.APIKey{}, err
    }
    return a, nil
}

func (s *Store) InsertAPIKey(ctx context.Context, a *data.APIKey) error {
    row := s.db.QueryRow(ctx, `INSERT INTO api_keys (consumer_id, key_hash, key_prefix, status, last_used_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`, a.ConsumerID, "", a.KeyPrefix, a.Status, a.LastUsedAt, a.ExpiresAt)
    return row.Scan(&a.ID, &a.CreatedAt)
}

func (s *Store) UpdateAPIKey(ctx context.Context, a *data.APIKey) error {
    _, err := s.db.Exec(ctx, `UPDATE api_keys SET key_prefix=$1, status=$2, last_used_at=$3, expires_at=$4 WHERE id=$5`, a.KeyPrefix, a.Status, a.LastUsedAt, a.ExpiresAt, a.ID)
    return err
}

func (s *Store) DeleteAPIKey(ctx context.Context, id string) error {
    _, err := s.db.Exec(ctx, `DELETE FROM api_keys WHERE id=$1`, id)
    return err
}

// Jobs CRUD
func (s *Store) GetJobByID(ctx context.Context, id string) (data.Job, error) {
    var j data.Job
    var payload []byte
    var result []byte
    row := s.db.QueryRow(ctx, `SELECT id, consumer_id, job_type, status, payload, result, error_message, started_at, completed_at, created_at FROM jobs WHERE id=$1`, id)
    if err := row.Scan(&j.ID, &j.ConsumerID, &j.JobType, &j.Status, &payload, &result, &j.ErrorMessage, &j.StartedAt, &j.CompletedAt, &j.CreatedAt); err != nil {
        return data.Job{}, err
    }
    j.Payload = json.RawMessage(payload)
    j.Result = json.RawMessage(result)
    return j, nil
}

func (s *Store) InsertJob(ctx context.Context, j *data.Job) error {
    row := s.db.QueryRow(ctx, `INSERT INTO jobs (consumer_id, job_type, payload, result, error_message, started_at, completed_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`, j.ConsumerID, j.JobType, j.Payload, j.Result, j.ErrorMessage, j.StartedAt, j.CompletedAt)
    return row.Scan(&j.ID, &j.CreatedAt)
}

func (s *Store) UpdateJob(ctx context.Context, j *data.Job) error {
    _, err := s.db.Exec(ctx, `UPDATE jobs SET job_type=$1, status=$2, payload=$3, result=$4, error_message=$5, started_at=$6, completed_at=$7 WHERE id=$8`, j.JobType, j.Status, j.Payload, j.Result, j.ErrorMessage, j.StartedAt, j.CompletedAt, j.ID)
    return err
}

func (s *Store) DeleteJob(ctx context.Context, id string) error {
    _, err := s.db.Exec(ctx, `DELETE FROM jobs WHERE id=$1`, id)
    return err
}

// Generate a consumer activity report for the given time range.
func (s *Store) GenerateReport(ctx context.Context, consumerID string, from, to time.Time) (*data.ConsumerActivityReport, error) {
    query := `
        SELECT
            c.id,
            c.name,
            COUNT(DISTINCT k.id) FILTER (WHERE k.status = 'active'),
            COUNT(DISTINCT k.id) FILTER (WHERE k.status = 'revoked'),
            COUNT(DISTINCT j.id) FILTER (WHERE j.status = 'queued'),
            COUNT(DISTINCT j.id) FILTER (WHERE j.status = 'processing'),
            COUNT(DISTINCT j.id) FILTER (WHERE j.status = 'completed'),
            COUNT(DISTINCT j.id) FILTER (WHERE j.status = 'failed')
        FROM consumers c
        LEFT JOIN api_keys k ON k.consumer_id = c.id
        LEFT JOIN jobs j ON j.consumer_id = c.id
            AND j.created_at >= $2
            AND j.created_at < $3
        WHERE c.id = $1
        GROUP BY c.id, c.name`

    report := &data.ConsumerActivityReport{
        From:        from,
        To:          to,
        GeneratedAt: time.Now(),
    }

    row := s.db.QueryRow(ctx, query, consumerID, from, to)
    err := row.Scan(
        &report.ConsumerID,
        &report.ConsumerName,
        &report.ActiveKeys,
        &report.RevokedKeys,
        &report.QueuedJobs,
        &report.ProcessingJobs,
        &report.CompletedJobs,
        &report.FailedJobs,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, sql.ErrNoRows
        }
        return nil, err
    }

    return report, nil
}
