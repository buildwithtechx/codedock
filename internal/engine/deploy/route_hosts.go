package deploy

import (
	"strings"

	"codedock/internal/models"
)

type RouteHost struct {
	Host       string
	PathPrefix string
	Protected  bool
}

func NormalizeRouteHosts(primary string, domains []models.DomainConfig) []RouteHost {
	return NormalizeRouteHostsWithProtection(primary, domains, nil)
}

func NormalizeRouteHostsWithProtection(primary string, domains []models.DomainConfig, isProtected func(string) bool) []RouteHost {
	hosts := make([]RouteHost, 0, len(domains)+1)
	seen := map[string]bool{}
	primary = strings.ToLower(strings.TrimSpace(primary))
	if primary != "" {
		seen[primary] = true
		prot := false
		if isProtected != nil {
			prot = isProtected(primary)
		}
		hosts = append(hosts, RouteHost{
			Host:       primary,
			PathPrefix: routePathPrefix(domainPathPrefix(primary, domains)),
			Protected:  prot,
		})
	}
	for _, domain := range domains {
		host := strings.ToLower(strings.TrimSpace(domain.DomainName))
		if host == "" || seen[host] {
			continue
		}
		seen[host] = true
		prot := false
		if isProtected != nil {
			prot = isProtected(host)
		}
		hosts = append(hosts, RouteHost{
			Host:       host,
			PathPrefix: routePathPrefix(domain.PathPrefix),
			Protected:  prot,
		})
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
