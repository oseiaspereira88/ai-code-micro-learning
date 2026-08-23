package eventstore

import (
	"errors"
	"fmt"
	"os"
)

// ErrLockHeld is returned by AcquireLock when another process already holds
// the workspace lock (requirement R7).
var ErrLockHeld = errors.New("eventstore: lock already held")

// Lock is an exclusive, filesystem-based lock for one workspace. Callers
// must release it exactly once, typically via defer, so cancellation or
// normal shutdown always frees it (requirement R7); a killed process still
// leaves the lock file behind, which is an inherent limit of file locks and
// is out of scope for this store (stale-lock recovery belongs to a
// dedicated operational spec).
type Lock struct {
	path string
	file *os.File
}

// AcquireLock creates path exclusively, failing with ErrLockHeld if it
// already exists.
func AcquireLock(path string) (*Lock, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return nil, ErrLockHeld
		}
		return nil, err
	}
	fmt.Fprintf(f, "%d\n", os.Getpid())
	return &Lock{path: path, file: f}, nil
}

// Release closes and removes the lock file. It is safe to call once; a
// second call returns an error rather than silently succeeding, so a bug
// double-releasing a lock is visible.
func (l *Lock) Release() error {
	if err := l.file.Close(); err != nil {
		return err
	}
	return os.Remove(l.path)
}
