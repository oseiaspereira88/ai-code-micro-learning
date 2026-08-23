package session

import (
	"github.com/oseiaspereira88/codinho/internal/curriculum"
	"github.com/oseiaspereira88/codinho/internal/learning"
)

// Window is the node a session's instructional window currently resolves
// to at a given depth. It never rewrites the canonical step tree
// (PROJECT.md §8.5, requirement R6; session-orchestration-disclosure
// Decision 1).
type Window struct {
	StepID string
	Kind   string
}

// deriveWindow walks challenge's authored tree to find the node at depth,
// following the first child at each level. It degrades gracefully: asking
// for a depth deeper than the content goes returns the deepest node
// actually authored, rather than failing.
func deriveWindow(ch curriculum.ChallengeAuthoring, depth learning.Depth) (Window, bool) {
	switch depth {
	case learning.DepthChallenge:
		return Window{StepID: ch.ID, Kind: "challenge"}, true
	case learning.DepthLayer:
		if len(ch.Layers) == 0 {
			return Window{}, false
		}
		return Window{StepID: ch.Layers[0].ID, Kind: "layer"}, true
	case learning.DepthMacro:
		step, ok := firstStep(ch)
		if !ok {
			return Window{}, false
		}
		return Window{StepID: step.ID, Kind: step.Kind}, true
	case learning.DepthMeso:
		step, ok := firstStep(ch)
		if !ok {
			return Window{}, false
		}
		target := descendTo(step, "meso")
		return Window{StepID: target.ID, Kind: target.Kind}, true
	case learning.DepthMicro:
		step, ok := firstStep(ch)
		if !ok {
			return Window{}, false
		}
		target := descendTo(step, "micro")
		return Window{StepID: target.ID, Kind: target.Kind}, true
	default:
		return Window{}, false
	}
}

// descendTo follows the first child at each level until it reaches a step
// of wantKind, or runs out of children — whichever comes first.
func descendTo(step curriculum.StepAuthoring, wantKind string) curriculum.StepAuthoring {
	current := step
	for current.Kind != wantKind && len(current.Children) > 0 {
		current = current.Children[0]
	}
	return current
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
