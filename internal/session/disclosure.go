package session

import "github.com/oseiaspereira88/ailearn/internal/learning"

// Disclosure is the client-facing view of a session's current revelation
// state (PROJECT.md §15.3 "disclosure").
type Disclosure struct {
	Level            string
	SolutionRevealed bool
}

// disclosureFor reports the wire-level name for policy's ceiling and
// whether the given step already had its solution revealed. Levels above
// "instruction" (pseudocode, fragment, solution) are only reached through
// hint_request, introduced by assistance-hints-detours; this server never
// discloses past level 0 on its own (requirement R3: never the answer key).
func disclosureFor(policy learning.DisclosurePolicy, step *learning.StepProgress) Disclosure {
	d := Disclosure{Level: "briefing"}
	if policy.MaxLevel >= learning.DisclosureGuidingQuestion {
		d.Level = "instruction"
	}
	if step != nil {
		d.SolutionRevealed = step.SolutionRevealed
	}
	return d
}
