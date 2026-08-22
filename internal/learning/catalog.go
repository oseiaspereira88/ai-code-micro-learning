// Package learning implements the pure pedagogical domain model: ontology,
// session and step invariants, policies, evidence and feedback. It never
// imports MCP, CLI, YAML or filesystem packages (PROJECT.md §12, §24.2) —
// loading, persisting and observing content belong to other specs.
package learning

// CatalogID identifies a curriculum entity (track, theme, challenge, ...)
// independently of its revision.
type CatalogID string

// ContentVersion pins the exact revision of a catalog entity that a session
// fixed at start; a session never silently observes a different version
// (PROJECT.md §12.3 invariant 1).
type ContentVersion string

// TrackID, ThemeID, ConceptID and CompetencyID identify catalog entities
// with stable IDs (PROJECT.md §7.1, requirement R1).
type (
	TrackID      CatalogID
	ThemeID      CatalogID
	ConceptID    CatalogID
	CompetencyID CatalogID
)

// Track is an ordered journey toward a learning goal.
type Track struct {
	ID      TrackID
	Version ContentVersion
	Themes  []ThemeID
}

// Theme groups challenges under a knowledge area.
type Theme struct {
	ID      ThemeID
	Version ContentVersion
}

// Concept is something a learner can understand. It is never graded by mere
// exposure; only competencies accrue evidence (PROJECT.md §7.2).
type Concept struct {
	ID      ConceptID
	Version ContentVersion
}

// Competency is something a learner can demonstrate through evidence.
type Competency struct {
	ID      CompetencyID
	Version ContentVersion
}
