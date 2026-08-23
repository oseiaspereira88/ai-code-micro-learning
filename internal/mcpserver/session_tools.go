package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/oseiaspereira88/ailearn/internal/application"
	"github.com/oseiaspereira88/ailearn/internal/learning"
)

type sessionStartArgs struct {
	ChallengeID string `json:"challenge_id" jsonschema:"ID of the challenge to fix for this session"`
	Mode        string `json:"mode,omitempty" jsonschema:"pedagogical mode: teaching, practice, review, debug, exploration or interview (default practice)"`
	RequestID   string `json:"request_id,omitempty" jsonschema:"idempotency key; retrying with the same value returns the original session"`
}

type sessionGetArgs struct {
	SessionID string `json:"session_id" jsonschema:"session ID returned by session_start"`
}

type instructionGetArgs struct {
	SessionID string `json:"session_id" jsonschema:"session ID returned by session_start"`
}

func registerSessionTools(server *mcp.Server, sessions *application.SessionService) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "session_start",
		Description: "Fix a challenge, mode and policies, and activate the session's first instructional step.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: false, IdempotentHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args sessionStartArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.ChallengeID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "challenge_id is required", false, nil), nil
		}
		result, err := sessions.Start(application.StartInput{
			ChallengeID: args.ChallengeID,
			Mode:        learning.PedagogicalMode(args.Mode),
			RequestID:   args.RequestID,
		})
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectSessionChanged, map[string]string{"objective": result.Objective})
		env.SessionID = string(result.SessionID)
		env.ActiveNode = &ActiveNode{ID: string(result.ActiveStep), Kind: "micro"}
		env.AllowedActions = []string{"instruction_get", "session_get"}
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "session_get",
		Description: "Return a session's state, active node and disclosure.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args sessionGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		result, err := sessions.Get(learning.SessionID(args.SessionID))
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]string{"state": string(result.State)})
		env.SessionID = string(result.SessionID)
		if result.ActiveStep != "" {
			env.ActiveNode = &ActiveNode{ID: string(result.ActiveStep), Kind: "micro"}
		}
		return nil, env, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "instruction_get",
		Description: "Deliver a single instruction at the authorized depth. Never returns future children or the solution.",
		Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true},
	}, func(_ context.Context, req *mcp.CallToolRequest, args instructionGetArgs) (*mcp.CallToolResult, Envelope, error) {
		requestID := requestIDFor(req)
		if args.SessionID == "" {
			return errorResult(), errorEnvelope(requestID, ErrCodeInvalidInput, "session_id is required", false, nil), nil
		}
		instruction, err := sessions.Instruction(learning.SessionID(args.SessionID))
		if err != nil {
			code, msg, retryable := mapError(err)
			return errorResult(), errorEnvelope(requestID, code, msg, retryable, nil), nil
		}
		env := okEnvelope(requestID, ProgressEffectNone, map[string]string{
			"objective": instruction.Objective,
			"scope":     instruction.Scope,
		})
		env.SessionID = args.SessionID
		env.ActiveNode = &ActiveNode{ID: string(instruction.StepID), Kind: "micro"}
		env.Disclosure = &Disclosure{Level: "instruction", SolutionRevealed: false}
		return nil, env, nil
	})
}
