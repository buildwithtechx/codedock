package ai

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestCuratedTools_IncludesDecisionTool(t *testing.T) {
	tools := curatedTools(nil, "general")
	if len(tools) == 0 {
		t.Fatalf("expected at least 1 curated tool, got 0")
	}

	found := false
	for _, tool := range tools {
		fn, ok := tool["function"].(map[string]any)
		if !ok {
			continue
		}
		if fn["name"] == "ask_user_decision" {
			found = true
			params, ok := fn["parameters"].(map[string]any)
			if !ok {
				t.Fatalf("expected parameters map")
			}
			props, ok := params["properties"].(map[string]any)
			if !ok {
				t.Fatalf("expected properties map")
			}
			if _, hasQuestions := props["questions"]; !hasQuestions {
				t.Errorf("expected ask_user_decision parameters to include 'questions'")
			}
			if _, hasTitle := props["title"]; !hasTitle {
				t.Errorf("expected ask_user_decision parameters to include 'title'")
			}
		}
	}

	if !found {
		t.Errorf("expected ask_user_decision tool in curatedTools")
	}
}

func TestCuratedTools_PersonaFiltering(t *testing.T) {
	if !IsReadOnlyPersona("logs") {
		t.Errorf("expected logs persona to be read-only")
	}
	if !IsReadOnlyPersona("security") {
		t.Errorf("expected security persona to be read-only")
	}
	if IsReadOnlyPersona("devops") {
		t.Errorf("expected devops persona to allow mutations")
	}
	if IsReadOnlyPersona("general") {
		t.Errorf("expected general persona to allow mutations")
	}

	if PersonaDisplayName("logs") != "Logs Investigator" {
		t.Errorf("unexpected display name: %s", PersonaDisplayName("logs"))
	}
	if PersonaDisplayName("security") != "Security Reviewer" {
		t.Errorf("unexpected display name: %s", PersonaDisplayName("security"))
	}
	if PersonaDisplayName("devops") != "DevOps Co-pilot" {
		t.Errorf("unexpected display name: %s", PersonaDisplayName("devops"))
	}
}

func TestBuildInstanceContext_Personas(t *testing.T) {
	ctx := context.Background()

	logsCtx := buildInstanceContext(ctx, "/projects", nil, Dependencies{}, "logs")
	if !strings.Contains(logsCtx, "Logs Investigator") || !strings.Contains(logsCtx, "strict read-only") {
		t.Errorf("expected logs investigator prompt with read-only constraint, got: %s", logsCtx)
	}

	secCtx := buildInstanceContext(ctx, "/projects", nil, Dependencies{}, "security")
	if !strings.Contains(secCtx, "Security Reviewer") || !strings.Contains(secCtx, "strict read-only") {
		t.Errorf("expected security reviewer prompt with read-only constraint, got: %s", secCtx)
	}

	devopsCtx := buildInstanceContext(ctx, "/projects", nil, Dependencies{}, "devops")
	if !strings.Contains(devopsCtx, "DevOps Co-pilot") {
		t.Errorf("expected devops prompt, got: %s", devopsCtx)
	}
}

func TestChatResponse_QuestionSerialization(t *testing.T) {
	resp := ChatResponse{
		Reply:      "Please confirm the following choice:",
		Provider:   "openai",
		Model:      "gpt-4o",
		Configured: true,
		Questions: []QuestionGroup{
			{
				ID:    "q_deploy_opts",
				Title: "Select Rollback Target",
				Questions: []QuestionItem{
					{
						ID:          "target_deployment",
						Title:       "Choose a previous deployment:",
						AllowOther:  false,
						MultiSelect: false,
						Options: []QuestionOption{
							{ID: "dep_1", Label: "v1.2.0 (Stable)", Description: "Deployed 2 hours ago"},
							{ID: "dep_2", Label: "v1.1.9 (Hotfix)", Description: "Deployed 1 day ago"},
						},
					},
				},
			},
		},
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal ChatResponse: %v", err)
	}

	var decoded ChatResponse
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		t.Fatalf("failed to unmarshal ChatResponse: %v", err)
	}

	if len(decoded.Questions) != 1 {
		t.Fatalf("expected 1 question group, got %d", len(decoded.Questions))
	}
	group := decoded.Questions[0]
	if group.Title != "Select Rollback Target" {
		t.Errorf("expected title 'Select Rollback Target', got %q", group.Title)
	}
	if len(group.Questions) != 1 {
		t.Fatalf("expected 1 question item, got %d", len(group.Questions))
	}
	item := group.Questions[0]
	if len(item.Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(item.Options))
	}
	if item.Options[0].Label != "v1.2.0 (Stable)" {
		t.Errorf("expected first option label 'v1.2.0 (Stable)', got %q", item.Options[0].Label)
	}
}
