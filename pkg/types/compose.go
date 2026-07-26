package types

type ComposeFile struct {
	Version  string                 `json:"version,omitempty" yaml:"version,omitempty"`
	Services map[string]interface{} `json:"services" yaml:"services"`
	Volumes  map[string]interface{} `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Networks map[string]interface{} `json:"networks,omitempty" yaml:"networks,omitempty"`
}
