package evidence

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestPutGetRoundTrip(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	id, err := s.Put([]byte("hello evidence"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := s.Get(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != "hello evidence" {
		t.Fatalf("got %q", got)
	}
}

func TestPutIsContentAddressedAndDeduplicates(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	id1, err := s.Put([]byte("same bytes"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	id2, err := s.Put([]byte("same bytes"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id1 != id2 {
		t.Fatalf("expected identical IDs for identical bytes, got %s and %s", id1, id2)
	}
}

func TestGetDetectsTampering(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	id, err := s.Put([]byte("original"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(dir, id)
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := os.WriteFile(path, []byte("tampered"), 0o600); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.Get(id)
	if !errors.Is(err, ErrIntegrityMismatch) {
		t.Fatalf("expected ErrIntegrityMismatch, got %v", err)
	}
}

func TestGetRejectsPathEscapingID(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = s.Get("../../etc/passwd")
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
}

func TestGetRejectsMalformedID(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cases := []string{"", "sha256-", "sha256-nothex", "sha1-" + "0000000000000000000000000000000000000000000000000000000000000000"}
	for _, id := range cases {
		if _, err := s.Get(id); !errors.Is(err, ErrInvalidID) {
			t.Fatalf("Get(%q): expected ErrInvalidID, got %v", id, err)
		}
	}
}
