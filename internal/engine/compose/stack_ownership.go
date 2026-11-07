package compose

import (
	"encoding/json"
	"fmt"
	"strings"
)

func validateStackOwnership(content, id string) error {
	var document map[string]any
	if err := json.Unmarshal([]byte(content), &document); err != nil {
		return err
	}
	for _, kind := range []string{"volumes", "networks"} {
		resources, _ := document[kind].(map[string]any)
		for name, value := range resources {
			settings, _ := value.(map[string]any)
			if custom, ok := settings["name"].(string); ok && custom != stackProject(id)+"_"+name {
				return fmt.Errorf("%s %s must use the stack-owned name; shared or externally named resources are unsupported", kind, name)
			}
			if kind == "networks" {
				if err := allowedComposeKeys(settings, "network "+name, "name external driver internal attachable enable_ipv4 enable_ipv6 ipam labels"); err != nil {
					return err
				}
				if driver, ok := settings["driver"].(string); ok && driver != "bridge" {
					return fmt.Errorf("network %s requires unsupported driver %s", name, driver)
				}
			}
		}
	}
	services, _ := document["services"].(map[string]any)
	for name, value := range services {
		service, _ := value.(map[string]any)
		if deployment, ok := service["deploy"].(map[string]any); ok {
			if err := allowedComposeKeys(deployment, name+".deploy", "replicas resources restart_policy"); err != nil {
				return err
			}
			if resources, ok := deployment["resources"].(map[string]any); ok {
				if err := allowedComposeKeys(resources, name+".resources", "limits reservations"); err != nil {
					return err
				}
				for _, kind := range []string{"limits", "reservations"} {
					if constraints, ok := resources[kind].(map[string]any); ok {
						if err := allowedComposeKeys(constraints, name+"."+kind, "cpus memory pids"); err != nil {
							return err
						}
					}
				}
			}
		}
		if labels, ok := service["labels"].(map[string]any); ok {
			for key := range labels {
				if strings.HasPrefix(key, "traefik.") || strings.HasPrefix(key, "codedock.") || strings.HasPrefix(key, "com.docker.compose.") {
					return fmt.Errorf("service %s cannot set control-plane routing or ownership labels", name)
				}
			}
		}
	}
	return nil
}
