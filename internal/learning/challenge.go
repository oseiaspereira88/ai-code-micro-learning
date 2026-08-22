package learning

// ChallengeID and ChallengeLayerID identify a challenge and one of its
// logical layers (compreensão, implementação, teste, ...). Layer is a lens,
// not a folder: a small challenge may have only a few (PROJECT.md §7.1).
type (
	ChallengeID      CatalogID
	ChallengeLayerID CatalogID
)

// StepID identifies a single node in a challenge's macro/meso/micro tree
// with a stable ID (requirement R1).
type StepID CatalogID

// StepKind distinguishes the three granularities a step can be authored or
// exposed at (PROJECT.md §7.1, §8.5).
type StepKind string

const (
	StepKindMacro StepKind = "macro"
	StepKindMeso  StepKind = "meso"
	StepKindMicro StepKind = "micro"
)

// Challenge is a complete problem referencing concepts and producing
// evidence toward competencies (PROJECT.md §7.1, §7.2).
type Challenge struct {
	ID      ChallengeID
	Version ContentVersion
	Layers  []ChallengeLayerID
	Root    StepID
}

// ChallengeLayer is a logical perspective relevant to a challenge.
type ChallengeLayer struct {
	ID ChallengeLayerID
}

// StepNode is one node of a challenge's macro/meso/micro tree. The tree
// structure lives in the flat Parent/Children references so the canonical
// tree never changes shape when a session narrows or widens its
// instructional window (PROJECT.md §8.5).
type StepNode struct {
	ID       StepID
	Kind     StepKind
	Parent   *StepID
	Children []StepID
}
