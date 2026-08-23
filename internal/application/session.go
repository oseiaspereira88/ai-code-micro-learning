package application

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/oseiaspereira88/ailearn/internal/curriculum"
	"github.com/oseiaspereira88/ailearn/internal/eventstore"
	"github.com/oseiaspereira88/ailearn/internal/learning"
)

// ErrSessionNotFound is returned when a session lookup finds nothing with
// the given ID.
var ErrSessionNotFound = errors.New("application: session not found")

// ErrChallengeHasNoSteps is a domain-adjacent error: the challenge exists
// but authored no macro step to start from.
var ErrChallengeHasNoSteps = errors.New("application: challenge has no steps")

type sessionRecord struct {
	session     *learning.LearningSession
	challengeID string
}

// SessionService starts and reads minimal learning sessions. It bridges the
// authored challenge shape (internal/curriculum) to the pure domain
// (internal/learning) for exactly the first step, and durably records
// session lifecycle events (internal/eventstore). Materializing a full
// step tree and progressive disclosure belongs to
// session-orchestration-disclosure (see mcp-stdio-foundation Decision 2).
type SessionService struct {
	catalog *CatalogService
	store   *eventstore.Store

	mu           sync.Mutex
	sessions     map[learning.SessionID]*sessionRecord
	startResults map[string]StartResult // keyed by client RequestID, for idempotent retries
	nextID       atomic.Uint64
}

// NewSessionService wires a SessionService to its catalog and event store.
func NewSessionService(catalog *CatalogService, store *eventstore.Store) *SessionService {
	return &SessionService{
		catalog:      catalog,
		store:        store,
		sessions:     map[learning.SessionID]*sessionRecord{},
		startResults: map[string]StartResult{},
	}
}

// StartInput is what a caller supplies to Start.
type StartInput struct {
	ChallengeID string
	Mode        learning.PedagogicalMode
	RequestID   string // idempotency key for retries (requirement R8, R4)
}

// StartResult is what Start returns on success.
type StartResult struct {
	SessionID  learning.SessionID
	ActiveStep learning.StepID
	Objective  string
}

// Start fixes challengeID, creates a new session at its first authored
// macro step, and durably records the session start. Calling Start again
// with the same RequestID returns the same result without creating a
// second session (requirement R4, R8: idempotent retry).
func (s *SessionService) Start(in StartInput) (StartResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if in.RequestID != "" {
		if cached, ok := s.startResults[in.RequestID]; ok {
			return cached, nil
		}
	}

	challenge, err := s.catalog.Challenge(in.ChallengeID)
	if err != nil {
		return StartResult{}, err
	}
	step, ok := firstStep(challenge)
	if !ok {
		return StartResult{}, ErrChallengeHasNoSteps
	}

	mode := in.Mode
	if mode == "" {
		mode = learning.ModePractice
	}
	policy, err := learning.NewSessionPolicy(mode, learning.DepthMicro, learning.HelpProgressive, learning.DisclosureGuidingQuestion, learning.EvaluationOnDemand, learning.AdvanceExplicit, nil)
	if err != nil {
		return StartResult{}, err
	}

	id := s.nextID.Add(1)
	sessionID := learning.SessionID(fmt.Sprintf("ses_%d", id))

	// Append first: the durable record of intent to start must exist
	// before the in-memory session is exposed to callers (local-event-store
	// R1/R3 ordering contract).
	if _, err := s.store.Append(string(sessionID), 0, in.RequestID, eventstore.EventSessionStarted, map[string]string{
		"challenge_id": in.ChallengeID,
		"mode":         string(mode),
	}); err != nil {
		return StartResult{}, err
	}

	session := learning.NewLearningSession(sessionID, policy, learning.CatalogRef{})
	if err := session.Transition(learning.SessionStateActive); err != nil {
		return StartResult{}, err
	}
	progress := learning.NewStepProgress(learning.StepID(step.ID))
	if err := progress.Transition(learning.StepStateAvailable, false); err != nil {
		return StartResult{}, err
	}
	if err := progress.Transition(learning.StepStateActive, false); err != nil {
		return StartResult{}, err
	}
	if err := session.SetActiveInstruction(progress); err != nil {
		return StartResult{}, err
	}

	s.sessions[sessionID] = &sessionRecord{session: session, challengeID: in.ChallengeID}

	result := StartResult{SessionID: sessionID, ActiveStep: progress.StepID, Objective: step.Instruction.Objective}
	if in.RequestID != "" {
		s.startResults[in.RequestID] = result
	}
	return result, nil
}

// GetResult is what Get returns.
type GetResult struct {
	SessionID  learning.SessionID
	State      learning.SessionState
	ActiveStep learning.StepID
}

// Get returns the current state of an in-memory session, or
// ErrSessionNotFound.
func (s *SessionService) Get(id learning.SessionID) (GetResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec, ok := s.sessions[id]
	if !ok {
		return GetResult{}, ErrSessionNotFound
	}
	var active learning.StepID
	if step := rec.session.ActiveStep(); step != nil {
		active = step.StepID
	}
	return GetResult{SessionID: id, State: rec.session.State, ActiveStep: active}, nil
}

// Instruction returns the single disclosed instruction for id's active
// step: objective and scope only, never constraints, children or a
// solution (PROJECT.md §15.6 "instruction_get").
type Instruction struct {
	StepID    learning.StepID
	Objective string
	Scope     string
}

// Instruction returns the active step's instruction for session id.
func (s *SessionService) Instruction(id learning.SessionID) (Instruction, error) {
	s.mu.Lock()
	rec, ok := s.sessions[id]
	s.mu.Unlock()
	if !ok {
		return Instruction{}, ErrSessionNotFound
	}
	active := rec.session.ActiveStep()
	if active == nil {
		return Instruction{}, ErrSessionNotFound
	}
	challenge, err := s.catalog.Challenge(rec.challengeID)
	if err != nil {
		return Instruction{}, err
	}
	step, ok := findStep(challenge, string(active.StepID))
	if !ok {
		return Instruction{}, ErrNotFound
	}
	return Instruction{StepID: active.StepID, Objective: step.Instruction.Objective, Scope: step.Instruction.Scope}, nil
}

func firstStep(ch curriculum.ChallengeAuthoring) (curriculum.StepAuthoring, bool) {
	for _, layer := range ch.Layers {
		if len(layer.MacroSteps) > 0 {
			return layer.MacroSteps[0], true
		}
	}
	return curriculum.StepAuthoring{}, false
}

func findStep(ch curriculum.ChallengeAuthoring, id string) (curriculum.StepAuthoring, bool) {
	for _, layer := range ch.Layers {
		if step, ok := searchSteps(layer.MacroSteps, id); ok {
			return step, true
		}
	}
	return curriculum.StepAuthoring{}, false
}

func searchSteps(steps []curriculum.StepAuthoring, id string) (curriculum.StepAuthoring, bool) {
	for _, step := range steps {
		if step.ID == id {
			return step, true
		}
		if found, ok := searchSteps(step.Children, id); ok {
			return found, true
		}
	}
	return curriculum.StepAuthoring{}, false
}
