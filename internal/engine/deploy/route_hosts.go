package deploy

import (
	"strings"

	"codedock/internal/models"
)

type RouteHost struct {
	Host       string
	PathPrefix string
}

func NormalizeRouteHosts(primary string, domains []models.DomainConfig) []RouteHost {
	hosts := make([]RouteHost, 0, len(domains)+1)
	seen := map[string]bool{}
	primary = strings.ToLower(strings.TrimSpace(primary))
	if primary != "" {
		seen[primary] = true
		hosts = append(hosts, RouteHost{Host: primary, PathPrefix: routePathPrefix(domainPathPrefix(primary, domains))})
	}
	for _, domain := range domains {
		host := strings.ToLower(strings.TrimSpace(domain.DomainName))
		if host == "" || seen[host] {
			continue
		}
		seen[host] = true
		hosts = append(hosts, RouteHost{Host: host, PathPrefix: routePathPrefix(domain.PathPrefix)})
	}
	return hosts
}

func domainPathPrefix(host string, domains []models.DomainConfig) string {
	for _, domain := range domains {
		if strings.ToLower(strings.TrimSpace(domain.DomainName)) == host {
			return domain.PathPrefix
		}
	}
	return ""
}

func routePathPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" || prefix == "/" {
		return ""
	}
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	return strings.TrimSuffix(prefix, "/")
}
