package ai

type AISettings struct {
	ID              string `json:"id"`
	DefaultProvider string `json:"defaultProvider"`
	OpenAIKey       string `json:"openAIKey,omitempty"`
	OpenAIModel     string `json:"openAIModel,omitempty"`
	AnthropicKey    string `json:"anthropicKey,omitempty"`
	AnthropicModel  string `json:"anthropicModel,omitempty"`
	GoogleKey       string `json:"googleKey,omitempty"`
	GoogleModel     string `json:"googleModel,omitempty"`
	MistralKey      string `json:"mistralKey,omitempty"`
	MistralModel    string `json:"mistralModel,omitempty"`
	GroqKey         string `json:"groqKey,omitempty"`
	GroqModel       string `json:"groqModel,omitempty"`
	DeepSeekKey     string `json:"deepSeekKey,omitempty"`
	DeepSeekModel   string `json:"deepSeekModel,omitempty"`
	XAIKey          string `json:"xaiKey,omitempty"`
	XAIModel        string `json:"xaiModel,omitempty"`
	MoonshotKey     string `json:"moonshotKey,omitempty"`
	MoonshotModel   string `json:"moonshotModel,omitempty"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type ResourceTag struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Subtitle string `json:"subtitle,omitempty"`
	Status   string `json:"status,omitempty"`
	Href     string `json:"href,omitempty"`
}

type ToolStep struct {
	ID     string `json:"id"`
	Tool   string `json:"tool"`
	Label  string `json:"label"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

type PendingAction struct {
	ActionID    string         `json:"actionId"`
	Tool        string         `json:"tool"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Method      string         `json:"method"`
	Path        string         `json:"path"`
	Arguments   map[string]any `json:"arguments,omitempty"`
	Destructive bool           `json:"destructive"`
}

type QuestionOption struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

type QuestionItem struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	AllowOther  bool             `json:"allowOther,omitempty"`
	MultiSelect bool             `json:"multiSelect,omitempty"`
	Options     []QuestionOption `json:"options,omitempty"`
}

type QuestionGroup struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	Questions []QuestionItem `json:"questions"`
}
