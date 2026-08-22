package learning

import (
	"errors"
	"testing"
)

func TestValidSessionTransitionTable(t *testing.T) {
	cases := []struct {
		from, to SessionState
		want     bool
	}{
		{SessionStateCreated, SessionStateActive, true},
		{SessionStateActive, SessionStateCompleted, true},
		{SessionStateActive, SessionStateAbandoned, true},
		{SessionStateActive, SessionStatePaused, true},
		{SessionStatePaused, SessionStateActive, true},
		{SessionStateCreated, SessionStateCompleted, false},
		{SessionStateCompleted, SessionStateActive, false},
		{SessionStateAbandoned, SessionStateActive, false},
		{SessionStatePaused, SessionStateCompleted, false},
	}
	for _, c := range cases {
		if got := ValidSessionTransition(c.from, c.to); got != c.want {
			t.Errorf("ValidSessionTransition(%s, %s) = %v, want %v", c.from, c.to, got, c.want)
		}
	}
}

func TestLearningSessionTransitionIsDeterministic(t *testing.T) {
	s := NewLearningSession("ses_1", testPolicy(t), CatalogRef{})
	if err := s.Transition(SessionStateActive); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.State != SessionStateActive {
		t.Fatalf("state = %s, want active", s.State)
	}
	if err := s.Transition(SessionStateCreated); !errors.Is(err, DomainError{Code: ErrCodeInvalidSessionTransition}) {
		t.Fatalf("expected ErrCodeInvalidSessionTransition, got %v", err)
	}
}

func TestSetActiveInstructionRejectsSecondActive(t *testing.T) {
	s := NewLearningSession("ses_1", testPolicy(t), CatalogRef{})
	first := NewStepProgress("step_1")
	second := NewStepProgress("step_2")

	if err := s.SetActiveInstruction(first); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err := s.SetActiveInstruction(second)
	if !errors.Is(err, DomainError{Code: ErrCodeInstructionAlreadyActive}) {
		t.Fatalf("expected ErrCodeInstructionAlreadyActive, got %v", err)
	}
	if s.ActiveStep().StepID != first.StepID {
		t.Fatalf("active step changed unexpectedly to %s", s.ActiveStep().StepID)
	}
}

func TestAdvanceRequiresCompletedStepUnlessOverride(t *testing.T) {
	s := NewLearningSession("ses_1", testPolicy(t), CatalogRef{})
	first := NewStepProgress("step_1")
	if err := s.SetActiveInstruction(first); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	next := NewStepProgress("step_2")
	err := s.Advance(next, false)
	if !errors.Is(err, DomainError{Code: ErrCodeAdvanceNotAllowed}) {
		t.Fatalf("expected ErrCodeAdvanceNotAllowed, got %v", err)
	}
	if s.ActiveStep().StepID != first.StepID {
		t.Fatalf("active step advanced despite rejected call")
	}

	if err := s.Advance(next, true); err != nil {
		t.Fatalf("override advance failed: %v", err)
	}
	if s.ActiveStep().StepID != next.StepID {
		t.Fatalf("active step = %s, want step_2", s.ActiveStep().StepID)
	}
}

func TestOpenDetourRequiresActiveInstruction(t *testing.T) {
	s := NewLearningSession("ses_1", testPolicy(t), CatalogRef{})
	if _, err := s.OpenDetour("why?"); !errors.Is(err, DomainError{Code: ErrCodeNoActiveInstruction}) {
		t.Fatalf("expected ErrCodeNoActiveInstruction, got %v", err)
	}
}

func TestOpenDetourKeepsOriginalStepActive(t *testing.T) {
	s := NewLearningSession("ses_1", testPolicy(t), CatalogRef{})
	active := NewStepProgress("step_1")
	if err := s.SetActiveInstruction(active); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	d, err := s.OpenDetour("why exported fields?")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.State != DetourStateOpened {
		t.Fatalf("detour state = %s, want opened", d.State)
	}
	if s.ActiveStep().StepID != active.StepID {
		t.Fatalf("active step changed during detour")
	}

	if err := d.Transition(DetourStateActive); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := d.Transition(DetourStateResolved); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := d.Transition(DetourStateReturned); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.ActiveStep().StepID != active.StepID {
		t.Fatalf("active step changed after detour returned")
	}
}

func TestValidStepTransitionTableAndOverrideSkip(t *testing.T) {
	cases := []struct {
		from, to StepState
		want     bool
	}{
		{StepStateLocked, StepStateAvailable, true},
		{StepStateAvailable, StepStateActive, true},
		{StepStateActive, StepStateEvaluated, true},
		{StepStateEvaluated, StepStateCompleted, true},
		{StepStateEvaluated, StepStateActive, true},
		{StepStateLocked, StepStateActive, false},
		{StepStateCompleted, StepStateActive, false},
	}
	for _, c := range cases {
		if got := ValidStepTransition(c.from, c.to); got != c.want {
			t.Errorf("ValidStepTransition(%s, %s) = %v, want %v", c.from, c.to, got, c.want)
		}
	}

	p := NewStepProgress("step_1")
	if err := p.Transition(StepStateSkipped, false); !errors.Is(err, DomainError{Code: ErrCodeInvalidStepTransition}) {
		t.Fatalf("expected rejection without override, got %v", err)
	}
	if err := p.Transition(StepStateSkipped, true); err != nil {
		t.Fatalf("override skip failed: %v", err)
	}
	if p.State != StepStateSkipped {
		t.Fatalf("state = %s, want skipped", p.State)
	}
}

func TestRevealSolutionInvalidatesAutonomyEvidence(t *testing.T) {
	a := Attempt{StepID: "step_1"}
	if !a.CountsAsAutonomyEvidence() {
		t.Fatal("fresh attempt should count as autonomy evidence")
	}
	a.SolutionRevealed = true
	if a.CountsAsAutonomyEvidence() {
		t.Fatal("attempt after solution reveal must not count as autonomy evidence")
	}
}

func testPolicy(t *testing.T) SessionPolicy {
	t.Helper()
	p, err := NewSessionPolicy(ModePractice, DepthMicro, HelpProgressive, DisclosureConceptOrAPI, EvaluationOnDemand, AdvanceExplicit, nil)
	if err != nil {
		t.Fatalf("unexpected error building policy: %v", err)
	}
	return p
}
