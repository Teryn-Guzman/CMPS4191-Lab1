package data

import (
	"encoding/json"
	"testing"
	"time"
)

func TestReportRequestJSON(t *testing.T) {
	input := []byte(`{"consumer_id":"consumer-1","from":"2026-01-01T00:00:00Z","to":"2027-01-01T00:00:00Z"}`)

	var request ReportRequest
	if err := json.Unmarshal(input, &request); err != nil {
		t.Fatalf("unmarshal report request: %v", err)
	}

	if request.ConsumerID != "consumer-1" {
		t.Fatalf("expected consumer-1, got %q", request.ConsumerID)
	}
	if !request.From.Equal(time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected from time: %s", request.From)
	}
	if !request.To.Equal(time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected to time: %s", request.To)
	}
}

func TestJobPayloadJSON(t *testing.T) {
	job := Job{Payload: json.RawMessage(`{"format":"csv"}`)}
	encoded, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("marshal job: %v", err)
	}
	if string(encoded) == "" {
		t.Fatal("expected encoded job JSON")
	}
}
