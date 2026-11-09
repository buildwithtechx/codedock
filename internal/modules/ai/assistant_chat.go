package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock/internal/modules/mcp"
)

func curatedTools(bridge *mcp.Bridge, persona string) []map[string]any {
	var out []map[string]any
	// Built-in interactive decision / question tool
	out = append(out, map[string]any{
		"type": "function",
		"function": map[string]any{
			"name":        "ask_user_decision",
			"description": "Ask the user one or more questions or prompt them for a structured decision before proceeding with a modification or choice.",
			"parameters": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title": map[string]any{
						"type":        "string",
						"description": "Short headline or title of the question group",
					},
					"questions": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"id":    map[string]any{"type": "string"},
								"title": map[string]any{"type": "string"},
								"allowOther": map[string]any{
									"type":        "boolean",
									"description": "Whether the user can provide a custom text response in addition to options",
								},
								"multiSelect": map[string]any{
									"type":        "boolean",
									"description": "Whether multiple options can be chosen",
								},
								"options": map[string]any{
									"type": "array",
									"items": map[string]any{
										"type": "object",
										"properties": map[string]any{
											"id":          map[string]any{"type": "string"},
											"label":       map[string]any{"type": "string"},
											"description": map[string]any{"type": "string"},
										},
										"required": []string{"id", "label"},
									},
								},
							},
							"required": []string{"id", "title"},
						},
					},
				},
				"required": []string{"title", "questions"},
			},
		},
	})

	if bridge == nil {
		return out
	}
	var names []string
	if IsReadOnlyPersona(persona) {
		names = []string{
			"get_projects", "get_projects_by_id", "get_apps", "get_apps_by_id",
			"get_databases", "get_databases_by_id", "get_deployments", "get_system_stats",
			"get_attention", "get_issues", "get_one_click", "get_backups", "get_jobs",
		}
	} else {
		names = []string{
			"get_projects", "get_projects_by_id", "get_apps", "get_apps_by_id",
			"post_apps_by_id_redeploy", "post_apps_by_id_restart", "post_apps_by_id_stop",
			"get_databases", "get_databases_by_id", "get_deployments", "get_system_stats",
			"get_attention", "get_issues", "get_one_click", "get_backups", "get_jobs",
		}
	}
	for _, n := range names {
		if t, ok := bridge.FindTool(n); ok {
			out = append(out, map[string]any{
				"type":     "function",
				"function": map[string]any{"name": t.Name, "description": t.Description, "parameters": t.InputSchema},
			})
		}
	}
	return out
}

func HandleChat(ctx context.Context, req ChatRequest, meta mcp.DispatchMeta, deps Dependencies) (ChatResponse, error) {
	sysPrompt := buildInstanceContext(ctx, req.CurrentPath, req.Mentions, deps, req.Persona)
	ep, err := resolveProviderEndpoint(deps.Settings, req.Provider, req.Model)
	if err != nil {
		reply, steps, pending := buildLocalFallback(ctx, req, deps)
		return ChatResponse{
			Reply: reply, Provider: "codedock-local", Model: "instance-copilot",
			Configured: false, Steps: steps, PendingActions: pending,
		}, nil
	}

	tools := curatedTools(deps.Bridge, req.Persona)
	msgs := append([]ChatMessage{}, req.Messages...)
	var steps []ToolStep
	var pending []PendingAction
	var executed []string
	var questions []QuestionGroup

	for hop := 0; hop < 3; hop++ {
		assistantMsg, callErr := callCompletion(ctx, ep, sysPrompt, msgs, tools)
		if callErr != nil {
			return ChatResponse{}, callErr
		}
		if len(assistantMsg.ToolCalls) == 0 {
			return ChatResponse{
				Reply: assistantMsg.Content, Provider: ep.Provider, Model: ep.Model,
				Configured: true, Steps: steps, PendingActions: pending, ExecutedActions: executed,
				Questions: questions,
			}, nil
		}

		msgs = append(msgs, assistantMsg)
		for _, tc := range assistantMsg.ToolCalls {
			if tc.Function.Name == "ask_user_decision" {
				var qGroup QuestionGroup
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &qGroup); err == nil {
					if qGroup.ID == "" {
						qGroup.ID = "q_" + uuid.NewString()[:8]
					}
					questions = append(questions, qGroup)
					steps = append(steps, ToolStep{
						ID:     tc.ID,
						Tool:   "ask_user_decision",
						Label:  "Waiting for your answer: " + qGroup.Title,
						Status: "pending",
					})
				}
				continue
			}

			toolDef, ok := deps.Bridge.FindTool(tc.Function.Name)
			if !ok {
				continue
			}
			var args map[string]any
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)

			if IsReadOnlyPersona(req.Persona) && !toolDef.Annotations.ReadOnlyHint {
				steps = append(steps, ToolStep{
					ID:     tc.ID,
					Tool:   toolDef.Name,
					Label:  fmt.Sprintf("Blocked: %s operates in read-only mode", PersonaDisplayName(req.Persona)),
					Status: "error",
					Detail: "Mutating actions are disallowed under this persona",
				})
				msgs = append(msgs, ChatMessage{
					Role:       "tool",
					ToolCallID: tc.ID,
					Content:    fmt.Sprintf(`{"error":"forbidden","message":"Persona %s is strictly read-only and cannot invoke mutating tool %s"}`, PersonaDisplayName(req.Persona), toolDef.Name),
				})
				continue
			}

			if !toolDef.Annotations.ReadOnlyHint && (req.Mode != "agent" || toolDef.Annotations.DestructiveHint) {
				actID := "act_" + uuid.NewString()[:8]
				act := PendingAction{
					ActionID: actID, Tool: toolDef.Name, Title: "Run " + toolDef.Name,
					Description: toolDef.Description, Method: toolDef.Method, Path: toolDef.Path,
					Arguments: args, Destructive: toolDef.Annotations.DestructiveHint,
				}
				pendingActionsMu.Lock()
				pendingActions[actID] = act
				pendingActionsMu.Unlock()
				pending = append(pending, act)
				steps = append(steps, ToolStep{ID: tc.ID, Tool: toolDef.Name, Label: "Awaiting approval: " + toolDef.Name, Status: "pending"})
				msgs = append(msgs, ChatMessage{
					Role: "tool", ToolCallID: tc.ID,
					Content: fmt.Sprintf(`{"status":"awaiting_user_confirmation","actionId":%q}`, actID),
				})
				continue
			}

			res := deps.Bridge.DispatchTool(ctx, toolDef, args, meta)
			status := "done"
			if !res.OK {
				status = "error"
			} else if !toolDef.Annotations.ReadOnlyHint {
				executed = append(executed, toolDef.Name)
			}
			steps = append(steps, ToolStep{ID: tc.ID, Tool: toolDef.Name, Label: toolDef.Description, Status: status})
			msgs = append(msgs, ChatMessage{Role: "tool", ToolCallID: tc.ID, Content: res.Text})
		}

		if len(questions) > 0 {
			replyText := assistantMsg.Content
			if strings.TrimSpace(replyText) == "" {
				replyText = "Please answer the question above to continue."
			}
			return ChatResponse{
				Reply: replyText, Provider: ep.Provider, Model: ep.Model,
				Configured: true, Steps: steps, PendingActions: pending, ExecutedActions: executed,
				Questions: questions,
			}, nil
		}
	}

	finalMsg, err := callCompletion(ctx, ep, sysPrompt, msgs, nil)
	if err != nil {
		return ChatResponse{}, err
	}
	return ChatResponse{
		Reply: finalMsg.Content, Provider: ep.Provider, Model: ep.Model,
		Configured: true, Steps: steps, PendingActions: pending, ExecutedActions: executed,
		Questions: questions,
	}, nil
}

func ConfirmAction(ctx context.Context, actionID string, meta mcp.DispatchMeta, deps Dependencies) (map[string]any, error) {
	pendingActionsMu.Lock()
	act, ok := pendingActions[actionID]
	if ok {
		delete(pendingActions, actionID)
	}
	pendingActionsMu.Unlock()
	if !ok {
		return nil, errors.New("action expired or already executed")
	}
	toolDef, exists := deps.Bridge.FindTool(act.Tool)
	if !exists {
		return nil, errors.New("tool not found")
	}
	res := deps.Bridge.DispatchTool(ctx, toolDef, act.Arguments, meta)
	return map[string]any{
		"actionId": act.ActionID, "tool": act.Tool, "ok": res.OK, "status": res.Status, "result": res.Data,
		"executedAt": time.Now().UTC().Format(time.RFC3339),
	}, nil
}
