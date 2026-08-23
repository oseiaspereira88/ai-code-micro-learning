package session

import (
	"testing"

	"github.com/oseiaspereira88/codinho/internal/learning"
)

func TestDisclosureForReportsBriefingBelowGuidingQuestion(t *testing.T) {
	d := disclosureFor(learning.DisclosurePolicy{MaxLevel: learning.DisclosureNone}, nil)
	if d.Level != "briefing" {
		t.Fatalf("level = %s, want briefing", d.Level)
	}
	if d.SolutionRevealed {
		t.Fatal("nil step must never report a revealed solution")
	}
}

func TestDisclosureForReportsInstructionAtOrAboveGuidingQuestion(t *testing.T) {
	d := disclosureFor(learning.DisclosurePolicy{MaxLevel: learning.DisclosureGuidingQuestion}, nil)
	if d.Level != "instruction" {
		t.Fatalf("level = %s, want instruction", d.Level)
	}
}

func TestDisclosureForReflectsStepSolutionRevealed(t *testing.T) {
	step := learning.NewStepProgress("step_1")
	if d := disclosureFor(learning.DisclosurePolicy{}, step); d.SolutionRevealed {
		t.Fatal("fresh step must not report a revealed solution")
	}
	step.RevealSolution()
	if d := disclosureFor(learning.DisclosurePolicy{}, step); !d.SolutionRevealed {
		t.Fatal("expected SolutionRevealed to reflect the step's own flag")
	}
}
