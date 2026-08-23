package application

import (
	"github.com/oseiaspereira88/codinho/internal/eventstore"
	"github.com/oseiaspereira88/codinho/internal/learning"
	"github.com/oseiaspereira88/codinho/internal/session"
)

// The session orchestration logic (lifecycle, disclosure, granularity)
// lives in internal/session; these aliases let mcpserver depend only on
// application's stable surface (session-orchestration-disclosure).
type (
	StartInput        = session.StartInput
	StartResult       = session.StartResult
	GetResult         = session.GetResult
	Instruction       = session.Instruction
	ConfigureInput    = session.ConfigureInput
	LifecycleResult   = session.LifecycleResult
	GranularityResult = session.GranularityResult
	Disclosure        = session.Disclosure
)

var (
	ErrSessionNotFound     = session.ErrSessionNotFound
	ErrChallengeHasNoSteps = session.ErrChallengeHasNoSteps
	ErrNoWindowAtDepth     = session.ErrNoWindowAtDepth
)

// SessionService adapts internal/session.Service for the MCP server.
type SessionService struct {
	svc *session.Service
}

// NewSessionService wires a SessionService to its catalog and event store.
func NewSessionService(catalog *CatalogService, store *eventstore.Store) *SessionService {
	return &SessionService{svc: session.New(catalog.catalog, store)}
}

// Start fixes a challenge and activates the session's first instructional
// step (requirement R1).
func (s *SessionService) Start(in StartInput) (StartResult, error) { return s.svc.Start(in) }

// Get returns a session's current state, active node, revision and
// disclosure (requirement R2).
func (s *SessionService) Get(id learning.SessionID) (GetResult, error) { return s.svc.Get(id) }

// Instruction returns the active step's disclosed instruction
// (requirement R3).
func (s *SessionService) Instruction(id learning.SessionID) (Instruction, error) {
	return s.svc.Instruction(id)
}

// Configure applies mutable session properties (requirement R4).
func (s *SessionService) Configure(in ConfigureInput) (LifecycleResult, error) {
	return s.svc.Configure(in)
}

// Pause, Resume and Finish control lifecycle without inferring step
// completion (requirement R5).
func (s *SessionService) Pause(id learning.SessionID, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.svc.Pause(id, expectedRevision, requestID)
}

func (s *SessionService) Resume(id learning.SessionID, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.svc.Resume(id, expectedRevision, requestID)
}

func (s *SessionService) Finish(id learning.SessionID, expectedRevision uint64, requestID string) (LifecycleResult, error) {
	return s.svc.Finish(id, expectedRevision, requestID)
}

// GranularityAdjust moves the session's instructional window without
// rewriting the canonical step tree (requirement R6).
func (s *SessionService) GranularityAdjust(id learning.SessionID, depth learning.Depth, expectedRevision uint64, requestID string) (GranularityResult, error) {
	return s.svc.GranularityAdjust(id, depth, expectedRevision, requestID)
}
