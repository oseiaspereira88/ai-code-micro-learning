package mcpserver

import (
	"errors"

	"github.com/oseiaspereira88/ailearn/internal/application"
	"github.com/oseiaspereira88/ailearn/internal/eventstore"
	"github.com/oseiaspereira88/ailearn/internal/learning"
)

// ErrorCode is one of the stable codes PROJECT.md §15.4 defines. Only the
// subset relevant to this minimal tool slice is mapped here; the remaining
// codes belong to the specs that introduce the operations that can raise
// them.
type ErrorCode string

const (
	ErrCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrCodeItemNotFound     ErrorCode = "ITEM_NOT_FOUND"
	ErrCodeSessionNotActive ErrorCode = "SESSION_NOT_ACTIVE"
	ErrCodeStateConflict    ErrorCode = "STATE_CONFLICT"
	ErrCodeInternalError    ErrorCode = "INTERNAL_ERROR"
)

// mapError translates an internal error into a stable error code, a safe
// message and whether the caller may retry (requirement R4). The message
// never includes the original error's text, a path or an environment
// value: those could leak internals (security requirement).
func mapError(err error) (ErrorCode, string, bool) {
	switch {
	case errors.Is(err, application.ErrNotFound):
		return ErrCodeItemNotFound, "the requested item does not exist in the catalog", false
	case errors.Is(err, application.ErrSessionNotFound):
		return ErrCodeSessionNotActive, "no active session with that ID", false
	case errors.Is(err, eventstore.ErrRevisionConflict):
		return ErrCodeStateConflict, "the session changed since your last read; fetch it again", true
	default:
		var domainErr learning.DomainError
		if errors.As(err, &domainErr) {
			return ErrCodeInvalidInput, "the request violates a domain rule: " + string(domainErr.Code), false
		}
		return ErrCodeInternalError, "an internal error occurred", true
	}
}
