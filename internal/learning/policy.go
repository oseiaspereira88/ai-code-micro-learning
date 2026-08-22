package learning

import "time"

// PedagogicalMode is one of the six independent session modes
// (PROJECT.md §8.2).
type PedagogicalMode string

const (
	ModeTeaching    PedagogicalMode = "teaching"
	ModePractice    PedagogicalMode = "practice"
	ModeReview      PedagogicalMode = "review"
	ModeDebug       PedagogicalMode = "debug"
	ModeExploration PedagogicalMode = "exploration"
	ModeInterview   PedagogicalMode = "interview"
)

func (m PedagogicalMode) valid() bool {
	switch m {
	case ModeTeaching, ModePractice, ModeReview, ModeDebug, ModeExploration, ModeInterview:
		return true
	}
	return false
}

// Depth is the initial instructional granularity a session is fixed at
// (PROJECT.md §8.3 item 2).
type Depth string

const (
	DepthChallenge Depth = "challenge"
	DepthLayer     Depth = "layer"
	DepthMacro     Depth = "macro"
	DepthMeso      Depth = "meso"
	DepthMicro     Depth = "micro"
)

func (d Depth) valid() bool {
	switch d {
	case DepthChallenge, DepthLayer, DepthMacro, DepthMeso, DepthMicro:
		return true
	}
	return false
}

// HelpPolicyKind is the assistance policy a session configures independently
// of mode or depth (PROJECT.md §8.3 item 3).
type HelpPolicyKind string

const (
	HelpFree             HelpPolicyKind = "free"
	HelpProgressive      HelpPolicyKind = "progressive"
	HelpLimited          HelpPolicyKind = "limited"
	HelpNoCode           HelpPolicyKind = "no_code"
	HelpNoHints          HelpPolicyKind = "no_hints"
	HelpOnlyAfterAttempt HelpPolicyKind = "only_after_attempt"
)

func (k HelpPolicyKind) valid() bool {
	switch k {
	case HelpFree, HelpProgressive, HelpLimited, HelpNoCode, HelpNoHints, HelpOnlyAfterAttempt:
		return true
	}
	return false
}

// HelpPolicy is the value object controlling how assistance is granted.
type HelpPolicy struct {
	Kind HelpPolicyKind
}

// DisclosureLevel is a rung of the assistance ladder, 0 (objective and
// scope only) through 6 (commented solution) (PROJECT.md §8.4).
type DisclosureLevel int

const (
	DisclosureNone            DisclosureLevel = 0
	DisclosureGuidingQuestion DisclosureLevel = 1
	DisclosureConceptOrAPI    DisclosureLevel = 2
	DisclosureLogicalOutline  DisclosureLevel = 3
	DisclosurePseudocode      DisclosureLevel = 4
	DisclosureSkeleton        DisclosureLevel = 5
	DisclosureSolution        DisclosureLevel = 6
)

func (l DisclosureLevel) valid() bool {
	return l >= DisclosureNone && l <= DisclosureSolution
}

// DisclosurePolicy caps how far the assistance ladder a session may climb.
type DisclosurePolicy struct {
	MaxLevel DisclosureLevel
}

// EvaluationPolicyKind is when evaluation may run (PROJECT.md §8.3 item 4).
type EvaluationPolicyKind string

const (
	EvaluationOnDemand       EvaluationPolicyKind = "on_demand"
	EvaluationOnStepComplete EvaluationPolicyKind = "on_step_complete"
	EvaluationOnlyAtEnd      EvaluationPolicyKind = "only_at_end"
)

func (k EvaluationPolicyKind) valid() bool {
	switch k {
	case EvaluationOnDemand, EvaluationOnStepComplete, EvaluationOnlyAtEnd:
		return true
	}
	return false
}

// AdvancePolicyKind is always explicit in the V1 (PROJECT.md §8.3 item 5).
type AdvancePolicyKind string

const AdvanceExplicit AdvancePolicyKind = "explicit"

func (k AdvancePolicyKind) valid() bool {
	return k == AdvanceExplicit
}

// SessionPolicy aggregates the seven independent dimensions a session
// configures (PROJECT.md §8.3). None of these dimensions is inferred from
// another.
type SessionPolicy struct {
	Mode         PedagogicalMode
	InitialDepth Depth
	Help         HelpPolicy
	Disclosure   DisclosurePolicy
	Evaluation   EvaluationPolicyKind
	Advance      AdvancePolicyKind
	TimeLimit    *time.Duration // nil means off
}

// NewSessionPolicy validates every dimension independently and fails at
// construction when any is invalid (non-functional requirement: value
// objects invalid at construction).
func NewSessionPolicy(mode PedagogicalMode, depth Depth, help HelpPolicyKind, disclosureMax DisclosureLevel, evaluation EvaluationPolicyKind, advance AdvancePolicyKind, timeLimit *time.Duration) (SessionPolicy, error) {
	if !mode.valid() {
		return SessionPolicy{}, newDomainError(ErrCodeInvalidValue, "pedagogical mode: "+string(mode))
	}
	if !depth.valid() {
		return SessionPolicy{}, newDomainError(ErrCodeInvalidValue, "depth: "+string(depth))
	}
	if !help.valid() {
		return SessionPolicy{}, newDomainError(ErrCodeInvalidValue, "help policy: "+string(help))
	}
	if !disclosureMax.valid() {
		return SessionPolicy{}, newDomainError(ErrCodeInvalidValue, "disclosure level out of range")
	}
	if !evaluation.valid() {
		return SessionPolicy{}, newDomainError(ErrCodeInvalidValue, "evaluation policy: "+string(evaluation))
	}
	if !advance.valid() {
		return SessionPolicy{}, newDomainError(ErrCodeInvalidValue, "advance policy: "+string(advance))
	}
	if timeLimit != nil && *timeLimit <= 0 {
		return SessionPolicy{}, newDomainError(ErrCodeInvalidValue, "time limit must be positive when set")
	}
	return SessionPolicy{
		Mode:         mode,
		InitialDepth: depth,
		Help:         HelpPolicy{Kind: help},
		Disclosure:   DisclosurePolicy{MaxLevel: disclosureMax},
		Evaluation:   evaluation,
		Advance:      advance,
		TimeLimit:    timeLimit,
	}, nil
}
