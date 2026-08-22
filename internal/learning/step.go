package learning

// StepState is the lifecycle of a single StepNode instance within a session
// (PROJECT.md §13.2).
type StepState string

const (
	StepStateLocked    StepState = "locked"
	StepStateAvailable StepState = "available"
	StepStateActive    StepState = "active"
	StepStateEvaluated StepState = "evaluated"
	StepStateCompleted StepState = "completed"
	// StepStateSkipped exists only via explicit override and never produces
	// domain evidence (PROJECT.md §13.2).
	StepStateSkipped StepState = "skipped"
)

// stepTransitions enumerates every non-override edge of PROJECT.md §13.2:
//
//	locked → available → active → evaluated → completed
//	                      ↑    ↘ active
//	                      └─────┘
var stepTransitions = map[StepState]map[StepState]bool{
	StepStateLocked:    {StepStateAvailable: true},
	StepStateAvailable: {StepStateActive: true},
	StepStateActive:    {StepStateEvaluated: true},
	StepStateEvaluated: {StepStateActive: true, StepStateCompleted: true},
	StepStateCompleted: {},
	StepStateSkipped:   {},
}

// ValidStepTransition reports whether moving a step from one state to
// another is allowed without an override.
func ValidStepTransition(from, to StepState) bool {
	return stepTransitions[from][to]
}

// StepProgress tracks one learner's progress through a single StepNode
// instance for the lifetime of a session.
type StepProgress struct {
	StepID           StepID
	State            StepState
	SolutionRevealed bool
}

// NewStepProgress starts a step at StepStateLocked.
func NewStepProgress(id StepID) *StepProgress {
	return &StepProgress{StepID: id, State: StepStateLocked}
}

// Transition moves the step to a new state, enforcing the state table.
// override bypasses the table only to reach StepStateSkipped, and is
// otherwise rejected (requirement R4).
func (p *StepProgress) Transition(to StepState, override bool) error {
	if to == StepStateSkipped {
		if !override {
			return newDomainError(ErrCodeInvalidStepTransition, "skipped requires an explicit override")
		}
		p.State = StepStateSkipped
		return nil
	}
	if !ValidStepTransition(p.State, to) {
		return newDomainError(ErrCodeInvalidStepTransition, string(p.State)+" -> "+string(to))
	}
	p.State = to
	return nil
}

// Complete marks the step completed. It requires the step to already be
// evaluated, unless override is set for an explicit exception (requirement
// R4; PROJECT.md §8.6 "Conclusão").
func (p *StepProgress) Complete(override bool) error {
	if p.State != StepStateEvaluated && !override {
		return newDomainError(ErrCodeStepNotCompletable, string(p.State))
	}
	p.State = StepStateCompleted
	return nil
}

// RevealSolution marks the step's solution as revealed. Per PROJECT.md §8.4,
// this does not itself change step state, and per invariant 8 it means any
// attempt already recorded, or recorded afterward, no longer counts as
// autonomy evidence (see Attempt.CountsAsAutonomyEvidence).
func (p *StepProgress) RevealSolution() {
	p.SolutionRevealed = true
}
