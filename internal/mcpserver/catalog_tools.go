package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/codinho/internal/application"
	"github.com/oseiaspereira88/codinho/internal/curriculum"
)

// catalogSearchArgs is catalog_search's strict input schema (requirement
// R5, R6).
type catalogSearchArgs struct {
	Kind  string `json:"kind,omitempty" jsonschema:"item kind to list: theme, concept, competency, track or challenge (default challenge)"`
	Theme string `json:"theme,omitempty" jsonschema:"optional theme ID to narrow results"`
}

// catalogGetArgs is catalog_get's strict input schema.
type catalogGetArgs struct {
	ID string `json:"id" jsonschema:"catalog item ID"`
}

func registerCatalogTools(server *mcp.Server, catalog *application.CatalogService) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "catalog_search",
		Description: "Search themes, concepts, competencies, tracks and challenges by kind and theme.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args catalogSearchArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		items := catalog.Search(application.SearchQuery{Kind: curriculum.ItemKind(args.Kind), Theme: args.Theme})
		return nil, okEnvelope(requestID, ProgressEffectNone, items), nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "catalog_get",
		Description: "Get one catalog item by ID and version, respecting the session's disclosure level.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args catalogGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.ID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "id is required", false, nil), nil
		}
		item, err := catalog.Get(args.ID)
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		return nil, okEnvelope(requestID, ProgressEffectNone, item), nil
	})
}
