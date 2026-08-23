package eventstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// ErrSnapshotMissing is returned by ReadSnapshot when no snapshot exists
// yet. Callers should rebuild the projection via Replay/ReplayAll instead
// (requirement R2).
var ErrSnapshotMissing = errors.New("eventstore: snapshot missing")

// ErrSnapshotInvalid is returned by ReadSnapshot when the file exists but is
// not well-formed JSON (requirement R2, R6).
var ErrSnapshotInvalid = errors.New("eventstore: snapshot invalid")

// WriteSnapshotAtomic replaces path with data via a same-directory
// temp-file-then-rename, so a reader never observes a partially written
// snapshot (requirement R3). Callers must call this only after the event it
// derives from has already been synced to the log (Store.Append fsyncs
// before returning), never before.
func WriteSnapshotAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".snapshot-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}
	if err := os.Chmod(tmpPath, 0o600); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, path)
}

// ReadSnapshot reads and validates the snapshot at path. A missing file is
// reported as ErrSnapshotMissing, not a bare os error, so callers can branch
// cleanly into a full replay.
func ReadSnapshot(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrSnapshotMissing
		}
		return nil, err
	}
	if !json.Valid(data) {
		return nil, ErrSnapshotInvalid
	}
	return data, nil
}
