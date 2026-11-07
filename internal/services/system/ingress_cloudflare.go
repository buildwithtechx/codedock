package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type IngressRule struct {
	Hostname      string         `json:"hostname,omitempty"`
	Service       string         `json:"service"`
	Path          string         `json:"path,omitempty"`
	OriginRequest map[string]any `json:"originRequest,omitempty"`
}

type TunnelConfig struct {
	Ingress []IngressRule `json:"ingress"`
}

type AccessApp struct {
	ID                string           `json:"id"`
	Name              string           `json:"name,omitempty"`
	Type              string           `json:"type,omitempty"`
	Domain            string           `json:"domain,omitempty"`
	AUD               string           `json:"aud,omitempty"`
	Policies          []string         `json:"policies,omitempty"`
	Destinations      []map[string]any `json:"destinations,omitempty"`
	SelfHostedDomains []string         `json:"self_hosted_domains,omitempty"`
}

type CloudflareIngressService struct {
	httpClient *http.Client
	baseURL    string
}

func NewCloudflareIngressService(client *http.Client) *CloudflareIngressService {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &CloudflareIngressService{
		httpClient: client,
		baseURL:    "https://api.cloudflare.com/client/v4",
	}
}

// Canonical generates a deterministic JSON representation sorting object keys
// and stripping empty or nil originRequest objects.
func Canonical(value any) string {
	switch v := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(v))
		for k, val := range v {
			if val == nil {
				continue
			}
			if subMap, ok := val.(map[string]any); ok && len(subMap) == 0 {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, len(keys))
		for i, k := range keys {
			parts[i] = fmt.Sprintf("%q:%s", k, Canonical(v[k]))
		}
		return "{" + strings.Join(parts, ",") + "}"
	case []any:
		parts := make([]string, len(v))
		for i, item := range v {
			parts[i] = Canonical(item)
		}
		return "[" + strings.Join(parts, ",") + "]"
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

// MergeIngress updates or inserts a rule before the final catch-all rule,
// enforcing strict conflict fences against manual modifications outside CodeDock.
func MergeIngress(cfg TunnelConfig, hostname string, prevRule, prevNextRule, nextRule *IngressRule) (TunnelConfig, error) {
	rules := cfg.Ingress
	if len(rules) == 0 {
		return cfg, fmt.Errorf("tunnel needs a final catch-all rule")
	}
	last := rules[len(rules)-1]
	if last.Hostname != "" || last.Path != "" || last.Service == "" {
		return cfg, fmt.Errorf("tunnel needs a final catch-all rule")
	}

	normHost := strings.ToLower(strings.TrimSpace(hostname))
	var matchIdx = -1
	for i, r := range rules[:len(rules)-1] {
		if strings.ToLower(r.Hostname) == normHost {
			if matchIdx >= 0 {
				return cfg, fmt.Errorf("multiple tunnel rules match %s", hostname)
			}
			matchIdx = i
		}
	}

	if matchIdx >= 0 {
		existing := rules[matchIdx]
		existingCanon := Canonical(existing)
		matchFound := false
		if prevRule != nil && Canonical(*prevRule) == existingCanon {
			matchFound = true
		}
		if prevNextRule != nil && Canonical(*prevNextRule) == existingCanon {
			matchFound = true
		}
		if !matchFound && prevRule != nil {
			return cfg, fmt.Errorf("the tunnel rule for %s changed outside CodeDock", hostname)
		}
	} else if nextRule != nil {
		for _, r := range rules {
			if strings.HasPrefix(r.Hostname, "*.") && strings.HasSuffix(normHost, r.Hostname[1:]) {
				return cfg, fmt.Errorf("a wildcard tunnel rule already covers %s", hostname)
			}
		}
	}

	copyRules := make([]IngressRule, 0, len(rules)+1)
	for i, r := range rules {
		if i == matchIdx {
			continue
		}
		if i == len(rules)-1 && nextRule != nil {
			copyRules = append(copyRules, *nextRule)
		}
		copyRules = append(copyRules, r)
	}

	return TunnelConfig{Ingress: copyRules}, nil
}

// ParseAccessApp validates Access apps, supporting domain-less launcher, SaaS,
// and private applications without requiring self-hosted domain or AUD fields (upstream #1302).
func ParseAccessApp(raw map[string]any) (AccessApp, error) {
	id, ok := raw["id"].(string)
	if !ok || id == "" {
		return AccessApp{}, fmt.Errorf("access app missing valid id")
	}
	app := AccessApp{ID: id}
	if name, ok := raw["name"].(string); ok {
		app.Name = name
	}
	if typ, ok := raw["type"].(string); ok {
		app.Type = typ
	}
	if dom, ok := raw["domain"].(string); ok {
		app.Domain = dom
	}
	if aud, ok := raw["aud"].(string); ok {
		app.AUD = aud
	}
	if rawPolicies, ok := raw["policies"].([]any); ok {
		for _, p := range rawPolicies {
			if s, ok := p.(string); ok {
				app.Policies = append(app.Policies, s)
			} else if m, ok := p.(map[string]any); ok {
				if pid, ok := m["id"].(string); ok {
					app.Policies = append(app.Policies, pid)
				}
			}
		}
	}
	return app, nil
}
