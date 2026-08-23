// Package application wires the pure domain (internal/learning), curriculum
// (internal/curriculum) and persistence (internal/eventstore,
// internal/evidence) packages into the use cases the MCP server exposes. It
// is the only layer that is allowed to know about all four.
package application

import (
	"errors"

	"github.com/oseiaspereira88/ailearn/internal/curriculum"
)

// ErrNotFound is returned when a catalog lookup finds nothing with the
// given ID.
var ErrNotFound = errors.New("application: not found")

// CatalogService adapts the immutable curriculum catalog for MCP tools.
type CatalogService struct {
	catalog *curriculum.Catalog
}

// NewCatalogService wraps an already-loaded, validated catalog.
func NewCatalogService(catalog *curriculum.Catalog) *CatalogService {
	return &CatalogService{catalog: catalog}
}

// SearchQuery narrows List by kind and theme (requirement R5, R6: "filtros
// básicos").
type SearchQuery struct {
	Kind  curriculum.ItemKind
	Theme string
}

// Search returns every sanitized item matching q. An empty Kind defaults to
// challenges, the slice's primary browsing surface.
func (s *CatalogService) Search(q SearchQuery) []curriculum.Item {
	kind := q.Kind
	if kind == "" {
		kind = curriculum.KindChallenge
	}
	return s.catalog.List(kind, q.Theme)
}

// Get returns the sanitized item with id, or ErrNotFound.
func (s *CatalogService) Get(id string) (curriculum.Item, error) {
	item, ok := s.catalog.Get(id)
	if !ok {
		return curriculum.Item{}, ErrNotFound
	}
	return item, nil
}

// Challenge returns the full, reserved-fields-cleared challenge record for
// id, or ErrNotFound.
func (s *CatalogService) Challenge(id string) (curriculum.ChallengeAuthoring, error) {
	ch, ok := s.catalog.Challenge(id)
	if !ok {
		return curriculum.ChallengeAuthoring{}, ErrNotFound
	}
	return ch, nil
}
