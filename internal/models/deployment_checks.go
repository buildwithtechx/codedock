package models

type GitHubDeploymentChecks struct {
	Enabled       bool     `json:"enabled"`
	Deployment    bool     `json:"deployment"`
	Services      []string `json:"services,omitempty"`
	IncludeErrors bool     `json:"includeErrors"`
}

func DefaultGitHubDeploymentChecks() GitHubDeploymentChecks {
	return GitHubDeploymentChecks{
		Enabled:       true,
		Deployment:    true,
		Services:      []string{"all"},
		IncludeErrors: true,
	}
}

func ResolveGitHubDeploymentChecks(custom *GitHubDeploymentChecks) GitHubDeploymentChecks {
	def := DefaultGitHubDeploymentChecks()
	if custom == nil {
		return def
	}
	res := *custom
	if len(res.Services) == 0 {
		res.Services = def.Services
	}
	return res
}

func (c *GitHubDeploymentChecks) TargetsService(serviceName string) bool {
	if !c.Enabled {
		return false
	}
	for _, s := range c.Services {
		if s == "all" || s == serviceName {
			return true
		}
	}
	return false
}
