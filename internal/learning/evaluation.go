package learning

// FeedbackType classifies consultative feedback (PROJECT.md §8.7). Feedback
// never mutates step state regardless of type (invariant 3) — Blocking only
// describes whether it should stop a completion request from *elsewhere*.
type FeedbackType string

const (
	FeedbackConfirmation FeedbackType = "confirmation"
	FeedbackQuestion     FeedbackType = "question"
	FeedbackExplanation  FeedbackType = "explanation"
	FeedbackRelation     FeedbackType = "relation"
	FeedbackSuggestion   FeedbackType = "suggestion"
	FeedbackIdiom        FeedbackType = "idiom"
	FeedbackRisk         FeedbackType = "risk"
	FeedbackViolation    FeedbackType = "violation"
	FeedbackError        FeedbackType = "error"
)

// blocksCompletionByDefault mirrors the "Bloqueia conclusão?" column of
// PROJECT.md §8.7. Idiom and Risk are context-dependent ("normalmente não" /
// "conforme critério") and default to non-blocking; a caller with more
// context can override via FeedbackRecord.Blocking.
func (t FeedbackType) blocksCompletionByDefault() bool {
	switch t {
	case FeedbackViolation, FeedbackError:
		return true
	default:
		return false
	}
}

func (t FeedbackType) valid() bool {
	switch t {
	case FeedbackConfirmation, FeedbackQuestion, FeedbackExplanation, FeedbackRelation,
		FeedbackSuggestion, FeedbackIdiom, FeedbackRisk, FeedbackViolation, FeedbackError:
		return true
	}
	return false
}

// FeedbackRecord is consultative: it describes observations, strengths,
// risks, suggestions and questions, approves nothing, consumes no attempt,
// and never completes or advances a step (PROJECT.md §8.6 invariant 3).
type FeedbackRecord struct {
	StepID   StepID
	Type     FeedbackType
	Text     string
	Blocking bool
}

// NewFeedbackRecord validates Type and applies the table default for
// Blocking unless overridden.
func NewFeedbackRecord(stepID StepID, feedbackType FeedbackType, text string, blockingOverride *bool) (FeedbackRecord, error) {
	if !feedbackType.valid() {
		return FeedbackRecord{}, newDomainError(ErrCodeInvalidValue, "feedback type: "+string(feedbackType))
	}
	if text == "" {
		return FeedbackRecord{}, newDomainError(ErrCodeInvalidValue, "feedback text is required")
	}
	blocking := feedbackType.blocksCompletionByDefault()
	if blockingOverride != nil {
		blocking = *blockingOverride
	}
	return FeedbackRecord{StepID: stepID, Type: feedbackType, Text: text, Blocking: blocking}, nil
}

// EvaluationVerdict is the outcome of comparing evidence with criteria
// (PROJECT.md §8.6 "Avaliação").
type EvaluationVerdict string

const (
	VerdictMet          EvaluationVerdict = "met"
	VerdictPartiallyMet EvaluationVerdict = "partially_met"
	VerdictNotMet       EvaluationVerdict = "not_met"
	VerdictUnverifiable EvaluationVerdict = "unverifiable"
)

// CriterionResult judges one evaluation criterion.
type CriterionResult struct {
	Name     string
	Blocking bool
	Verdict  EvaluationVerdict
}

// Evaluation compares evidence with criteria. It never completes or
// advances a step by itself (PROJECT.md §12.3 invariant 4); those are
// separate, explicit operations on StepProgress/LearningSession.
type Evaluation struct {
	StepID    StepID
	Criteria  []CriterionResult
	AttemptID *EvidenceID // set only when the learner deliberately submitted
}

// HasBlockingFailure reports whether any blocking criterion was not met.
func (e Evaluation) HasBlockingFailure() bool {
	for _, c := range e.Criteria {
		if c.Blocking && c.Verdict != VerdictMet {
			return true
		}
	}
	return false
}
