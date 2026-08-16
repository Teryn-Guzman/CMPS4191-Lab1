package data

import (
	"encoding/json"
	"time"
)

type Consumer struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Status    string    `json:"status"`
    Version   int       `json:"version"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type APIKey struct {
    ID         string     `json:"id"`
    ConsumerID string     `json:"consumer_id"`
    KeyPrefix  string     `json:"key_prefix"`
    Status     string     `json:"status"`
    LastUsedAt *time.Time `json:"last_used_at,omitempty"`
    ExpiresAt  *time.Time `json:"expires_at,omitempty"`
    CreatedAt  time.Time  `json:"created_at"`
}

type Job struct {
    ID           string          `json:"id"`
    ConsumerID   string          `json:"consumer_id"`
    JobType      string          `json:"job_type"`
    Status       string          `json:"status"`
    Payload      json.RawMessage `json:"payload"`
    Result       json.RawMessage `json:"result,omitempty"`
    ErrorMessage *string         `json:"error_message,omitempty"`
    StartedAt    *time.Time      `json:"started_at,omitempty"`
    CompletedAt  *time.Time      `json:"completed_at,omitempty"`
    CreatedAt    time.Time       `json:"created_at"`
}
