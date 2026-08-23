// Package evidence stores evidence blobs content-addressed by SHA-256, so
// large output never lives in the event log itself (PROJECT.md §18.1, §18.2)
// and every read is verified against tampering (requirement R5).
package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ErrIntegrityMismatch is returned by Get when the stored bytes no longer
// hash to the requested ID.
var ErrIntegrityMismatch = errors.New("evidence: integrity mismatch")

// ErrInvalidID is returned when an ID is not a well-formed "sha256-<hex>"
// identifier. Rejecting it before any path join keeps a caller-supplied ID
// from ever escaping the evidence directory (security requirement).
var ErrInvalidID = errors.New("evidence: invalid id")

const idPrefix = "sha256-"

// Store persists evidence blobs under one directory.
type Store struct {
	dir string
}

// Open ensures dir exists and returns a Store rooted at it.
func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

// Put stores data and returns its content-addressed ID. Storing the same
// bytes twice is a cheap no-op and returns the same ID (PROJECT.md §18.1).
func (s *Store) Put(data []byte) (string, error) {
	id := idFor(data)
	path := filepath.Join(s.dir, id)
	if _, err := os.Stat(path); err == nil {
		return id, nil
	}
	tmp, err := os.CreateTemp(s.dir, ".evidence-*.tmp")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return "", err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return "", err
	}
	if err := os.Chmod(tmpPath, 0o400); err != nil { // evidence is immutable once written
		os.Remove(tmpPath)
		return "", err
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return "", err
	}
	return id, nil
}

// Get reads the evidence identified by id and verifies it still hashes to
// id before returning it (requirement R5).
func (s *Store) Get(id string) ([]byte, error) {
	if !validID(id) {
		return nil, ErrInvalidID
	}
	data, err := os.ReadFile(filepath.Join(s.dir, id))
	if err != nil {
		return nil, err
	}
	if idFor(data) != id {
		return nil, ErrIntegrityMismatch
	}
	return data, nil
}

func idFor(data []byte) string {
	sum := sha256.Sum256(data)
	return idPrefix + hex.EncodeToString(sum[:])
}

func validID(id string) bool {
	digest := strings.TrimPrefix(id, idPrefix)
	if len(digest) != sha256.Size*2 || digest == id {
		return false
	}
	for _, c := range digest {
		isDigit := c >= '0' && c <= '9'
		isLower := c >= 'a' && c <= 'f'
		if !isDigit && !isLower {
			return false
		}
	}
	return true
}
