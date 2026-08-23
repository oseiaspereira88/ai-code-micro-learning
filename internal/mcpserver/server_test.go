package mcpserver

import (
	"errors"
	"testing"

	"github.com/oseiaspereira88/ailearn/internal/application"
	"github.com/oseiaspereira88/ailearn/internal/eventstore"
	"github.com/oseiaspereira88/ailearn/internal/learning"
)

func TestCoreInstructionsFitWithinFirst512Chars(t *testing.T) {
	if len(coreInstructions) > 512 {
		t.Fatalf("coreInstructions is %d chars, must fit within 512", len(coreInstructions))
	}
	if len(Instructions) < len(coreInstructions) {
		t.Fatal("Instructions must start with coreInstructions")
	}
	if Instructions[:len(coreInstructions)] != coreInstructions {
		t.Fatal("coreInstructions must be a prefix of Instructions")
	}
}

func TestRequestIDForIsUnique(t *testing.T) {
	a := requestIDFor(nil)
	b := requestIDFor(nil)
	if a == b {
		t.Fatalf("expected distinct request IDs, got %q twice", a)
	}
}

func TestErrorResultIsMarkedAsError(t *testing.T) {
	r := errorResult()
	if !r.IsError {
		t.Fatal("errorResult() must set IsError")
	}
}

func TestMapErrorTranslatesKnownErrors(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want ErrorCode
	}{
		{"not found", application.ErrNotFound, ErrCodeItemNotFound},
		{"session not found", application.ErrSessionNotFound, ErrCodeSessionNotActive},
		{"revision conflict", eventstore.ErrRevisionConflict, ErrCodeStateConflict},
		{"domain error", learning.DomainError{Code: learning.ErrCodeInvalidValue}, ErrCodeInvalidInput},
		{"challenge has no steps", application.ErrChallengeHasNoSteps, ErrCodeInvalidInput},
		{"no window at depth", application.ErrNoWindowAtDepth, ErrCodeInvalidInput},
		{"unknown", errors.New("boom"), ErrCodeInternalError},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, msg, _ := mapError(c.err)
			if code != c.want {
				t.Fatalf("mapError(%v) code = %s, want %s", c.err, code, c.want)
			}
			if msg == "" {
				t.Fatal("mapError must always return a non-empty message")
			}
		})
	}
}

func TestMapErrorNeverLeaksOriginalMessage(t *testing.T) {
	sensitive := errors.New("open /home/user/.ssh/id_rsa: permission denied")
	_, msg, _ := mapError(sensitive)
	if msg == sensitive.Error() {
		t.Fatal("mapError must not pass internal error text straight through to the client")
	}
}
