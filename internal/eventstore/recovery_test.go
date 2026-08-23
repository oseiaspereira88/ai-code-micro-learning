package eventstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func writeLog(t *testing.T, lines ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "events.jsonl")
	content := ""
	for _, l := range lines {
		content += l + "\n"
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return path
}

func encodeEvent(t *testing.T, streamID string, revision uint64, eventType EventType, schemaVersion int, payload string) string {
	t.Helper()
	ev := Event{
		SchemaVersion: schemaVersion,
		ID:            streamID + "@" + string(eventType),
		StreamID:      streamID,
		Revision:      revision,
		Type:          eventType,
		RecordedAt:    "2026-08-22T00:00:00Z",
		Payload:       json.RawMessage(payload),
	}
	raw, err := json.Marshal(ev)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return string(raw)
}

func TestReadEventsReturnsEmptyResultForMissingFile(t *testing.T) {
	result := ReadEvents(filepath.Join(t.TempDir(), "does-not-exist.jsonl"), nil)
	if result.Err != nil || result.Truncated || len(result.Events) != 0 {
		t.Fatalf("unexpected result for missing file: %+v", result)
	}
}

func TestReadEventsToleratesTruncatedFinalLine(t *testing.T) {
	good := encodeEvent(t, "ses_1", 1, EventSessionStarted, CurrentEventSchemaVersion, `{}`)
	path := writeLog(t, good)
	// Append a truncated line directly (writeLog would add a clean newline,
	// so append the broken fragment separately without one).
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := f.WriteString(`{"schema_version":1,"id":"ses_1@2","stream_id":"ses_1","revi`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	f.Close()

	result := ReadEvents(path, nil)
	if result.Err != nil {
		t.Fatalf("a truncated final line must not be reported as an error: %v", result.Err)
	}
	if !result.Truncated {
		t.Fatal("expected Truncated = true")
	}
	if len(result.Events) != 1 {
		t.Fatalf("valid prefix lost: got %d events, want 1", len(result.Events))
	}
}

func TestReadEventsReportsCorruptedMiddleLineWithoutLosingPrefix(t *testing.T) {
	good1 := encodeEvent(t, "ses_1", 1, EventSessionStarted, CurrentEventSchemaVersion, `{}`)
	good2 := encodeEvent(t, "ses_1", 3, EventSessionFinished, CurrentEventSchemaVersion, `{}`)
	path := writeLog(t, good1, "not json at all {{{", good2)

	result := ReadEvents(path, nil)
	if !errors.Is(result.Err, ErrCorruptedLog) {
		t.Fatalf("expected ErrCorruptedLog, got %v", result.Err)
	}
	if len(result.Events) != 1 {
		t.Fatalf("corrupted middle line must preserve the valid prefix: got %d events, want 1", len(result.Events))
	}
}

func TestReadEventsRejectsIncompatibleVersionWithoutUpcaster(t *testing.T) {
	future := encodeEvent(t, "ses_1", 1, EventSessionStarted, CurrentEventSchemaVersion+1, `{}`)
	path := writeLog(t, future)

	result := ReadEvents(path, nil)
	if !errors.Is(result.Err, ErrIncompatibleEventVersion) {
		t.Fatalf("expected ErrIncompatibleEventVersion, got %v", result.Err)
	}
}

func TestReadEventsUpcastsRegisteredVersion(t *testing.T) {
	oldLine := encodeEvent(t, "ses_1", 1, EventSessionStarted, 0, `{"mode":"legacy"}`)
	path := writeLog(t, oldLine)

	registry := NewUpcasterRegistry()
	registry.Register(EventSessionStarted, 0, func(raw json.RawMessage) (json.RawMessage, error) {
		return json.RawMessage(`{"mode":"legacy","upcasted":true}`), nil
	})

	result := ReadEvents(path, registry)
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if len(result.Events) != 1 {
		t.Fatalf("got %d events, want 1", len(result.Events))
	}
	if result.Events[0].SchemaVersion != CurrentEventSchemaVersion {
		t.Fatalf("schema version = %d, want %d", result.Events[0].SchemaVersion, CurrentEventSchemaVersion)
	}
	if string(result.Events[0].Payload) != `{"mode":"legacy","upcasted":true}` {
		t.Fatalf("payload not upcasted: %s", result.Events[0].Payload)
	}
}
