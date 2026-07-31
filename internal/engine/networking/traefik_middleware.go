package networking

import (
	"encoding/json"
	"fmt"
	"strings"

	"codedock.run/codedock/internal/models"
)

func BuildMiddlewareLabels(serviceName string, rules []*models.RouteRule) map[string]string {
	labels := make(map[string]string)
	var middlewareNames []string

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		key := sanitizeLabelKey(serviceName + "-" + rule.Name)
		var ruleLabels map[string]string

		switch rule.RuleType {
		case models.RouteRuleTypeRateLimit:
			ruleLabels = buildRateLimitLabels(key, rule.SpecJSON)
			if len(ruleLabels) > 0 {
				middlewareNames = append(middlewareNames, key+"-ratelimit@docker")
			}
		case models.RouteRuleTypeIPAllowlist:
			ruleLabels = buildIPAllowlistLabels(key, rule.SpecJSON)
			if len(ruleLabels) > 0 {
				middlewareNames = append(middlewareNames, key+"-ipallow@docker")
			}
		case models.RouteRuleTypeHeaders:
			ruleLabels = buildHeadersLabels(key, rule.SpecJSON)
			if len(ruleLabels) > 0 {
				middlewareNames = append(middlewareNames, key+"-headers@docker")
			}
		}

		for k, v := range ruleLabels {
			labels[k] = v
		}
	}

	if len(middlewareNames) > 0 {
		svcKey := sanitizeLabelKey(serviceName)
		labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", svcKey)] = strings.Join(middlewareNames, ",")
	}

	return labels
}

func buildRateLimitLabels(key, specJSON string) map[string]string {
	var spec models.RateLimitSpec
	if err := json.Unmarshal([]byte(specJSON), &spec); err != nil || spec.Average <= 0 {
		return nil
	}
	period := spec.Period
	if period == "" {
		period = "1s"
	}
	prefix := fmt.Sprintf("traefik.http.middlewares.%s-ratelimit.ratelimit", key)
	return map[string]string{
		prefix + ".average": fmt.Sprintf("%d", spec.Average),
		prefix + ".burst":   fmt.Sprintf("%d", spec.Burst),
		prefix + ".period":  period,
	}
}

func buildIPAllowlistLabels(key, specJSON string) map[string]string {
	var spec models.IPListSpec
	if err := json.Unmarshal([]byte(specJSON), &spec); err != nil || len(spec.CIDRs) == 0 {
		return nil
	}
	return map[string]string{
		fmt.Sprintf("traefik.http.middlewares.%s-ipallow.ipallowlist.sourcerange", key): strings.Join(spec.CIDRs, ","),
	}
}

func buildHeadersLabels(key, specJSON string) map[string]string {
	var spec models.HeadersSpec
	if err := json.Unmarshal([]byte(specJSON), &spec); err != nil {
		return nil
	}
	labels := make(map[string]string)
	prefix := fmt.Sprintf("traefik.http.middlewares.%s-headers.headers", key)
	for k, v := range spec.Set {
		labels[fmt.Sprintf("%s.customrequestheaders.%s", prefix, k)] = v
	}
	for _, k := range spec.Remove {
		labels[fmt.Sprintf("%s.customrequestheaders.%s", prefix, k)] = ""
	}
	return labels
}

func ApplyMiddlewareLabels(existing map[string]string, serviceName string, rules []*models.RouteRule) map[string]string {
	extra := BuildMiddlewareLabels(serviceName, rules)
	if existing == nil {
		existing = make(map[string]string)
	}
	routerKey := fmt.Sprintf("traefik.http.routers.%s.middlewares", sanitizeLabelKey(serviceName))
	for k, v := range extra {
		if k == routerKey {
			if prev, ok := existing[routerKey]; ok && prev != "" {
				existing[routerKey] = prev + "," + v
			} else {
				existing[routerKey] = v
			}
		} else {
			existing[k] = v
		}
	}
	return existing
}

func sanitizeLabelKey(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}
