package compose

import (
	"codedock.run/codedock/internal/models"
	"fmt"
	"gopkg.in/yaml.v3"
	"net/url"
	"path"
	"strings"
)

func ResolveRepositoryContexts(request models.ComposeStackRequest) (string, error) {
	if request.RepositoryURL == "" {
		return request.Content, nil
	}
	if err := validateRemoteContext(request.RepositoryURL); err != nil {
		return "", err
	}
	parsed, err := url.Parse(request.RepositoryURL)
	if err != nil {
		return "", err
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("set the repository branch and root directory separately from its URL")
	}
	if strings.ContainsAny(request.Branch, " :#?\t\r\n") {
		return "", fmt.Errorf("invalid repository branch")
	}
	root := strings.TrimPrefix(strings.ReplaceAll(request.RootDirectory, "\\", "/"), "/")
	var document map[string]any
	if err := yaml.Unmarshal([]byte(request.Content), &document); err != nil {
		return "", err
	}
	services, _ := document["services"].(map[string]any)
	for name, value := range services {
		service, _ := value.(map[string]any)
		var context string
		settings, expanded := service["build"].(map[string]any)
		if expanded {
			context, _ = settings["context"].(string)
			if context == "" {
				context = "."
			}
		} else {
			context, _ = service["build"].(string)
		}
		if context == "" || strings.Contains(context, "://") {
			continue
		}
		context = strings.ReplaceAll(context, "\\", "/")
		if strings.HasPrefix(context, "/") || strings.Contains(context, ":") {
			return "", fmt.Errorf("service %s requires a relative repository build context", name)
		}
		directory := path.Clean(path.Join(root, context))
		if directory == ".." || strings.HasPrefix(directory, "../") {
			return "", fmt.Errorf("service %s build context escapes its repository", name)
		}
		if directory == "." {
			directory = ""
		}
		source := *parsed
		source.Fragment = request.Branch + ":" + directory
		if expanded {
			settings["context"] = source.String()
		} else {
			service["build"] = source.String()
		}
	}
	output, err := yaml.Marshal(document)
	if err != nil {
		return "", err
	}
	return string(output), nil
}
