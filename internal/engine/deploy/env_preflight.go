package deploy

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"codedock/internal/engine/build"
	"codedock/internal/models"
)

// Regex to detect Docker Compose required variable syntax:
// ${VAR:?error} or ${VAR?error} or ${VAR:? }
var requiredVarPattern = regexp.MustCompile(`\$\{([A-Za-z0-9_]+):?\?([^}]*)\}`)

// Regex to detect basic bash/compose interpolation ${VAR} or $VAR
var basicVarPattern = regexp.MustCompile(`\$\{([A-Za-z0-9_]+)\}`)

// MissingEnvVariable records an unsatisfied mandatory environment variable.
type MissingEnvVariable struct {
	ServiceName string
	Variable    string
	Message     string
}

// FindMissingRequiredEnv scans values for required compose syntax ${VAR:?msg} or ${VAR?msg}
// and checks if the referenced variable exists in the resolved environment map.
func FindMissingRequiredEnv(serviceName string, rawValues map[string]string, resolvedEnv map[string]string) []MissingEnvVariable {
	var missing []MissingEnvVariable
	seen := make(map[string]bool)

	for _, raw := range rawValues {
		matches := requiredVarPattern.FindAllStringSubmatch(raw, -1)
		for _, m := range matches {
			if len(m) < 2 {
				continue
			}
			varName := m[1]
			customMsg := ""
			if len(m) >= 3 {
				customMsg = strings.TrimSpace(m[2])
			}

			val, exists := resolvedEnv[varName]
			isColonQuery := strings.Contains(m[0], ":?")
			// If :? then empty string is also considered missing; if ? then only unset is missing
			if !exists || (isColonQuery && strings.TrimSpace(val) == "") {
				key := serviceName + ":" + varName
				if !seen[key] {
					seen[key] = true
					missing = append(missing, MissingEnvVariable{
						ServiceName: serviceName,
						Variable:    varName,
						Message:     customMsg,
					})
				}
			}
		}
	}

	return missing
}

// FormatRequiredEnvError formats missing required variables matching OpenShip's standard error format.
func FormatRequiredEnvError(missing []MissingEnvVariable) error {
	if len(missing) == 0 {
		return nil
	}

	byService := make(map[string][]string)
	for _, m := range missing {
		byService[m.ServiceName] = append(byService[m.ServiceName], m.Variable)
	}

	var serviceNames []string
	for name := range byService {
		serviceNames = append(serviceNames, name)
	}
	sort.Strings(serviceNames)

	var failures []string
	for _, name := range serviceNames {
		vars := byService[name]
		sort.Strings(vars)
		failures = append(failures, fmt.Sprintf(`Service "%s": %s`, name, strings.Join(vars, ", ")))
	}

	return fmt.Errorf("Required Compose environment variables are not configured. %s", strings.Join(failures, "; "))
}

// AssertRequiredServiceEnvironment performs an early pre-flight validation on the service's
// environment before git cloning or docker image builds begin.
func AssertRequiredServiceEnvironment(ctx context.Context, deployer *Deployer, app *models.AppService) error {
	if deployer == nil || app == nil {
		return nil
	}

	envVarsMap, err := deployer.getEnvironmentVariables(ctx, app, nil)
	if err != nil {
		return fmt.Errorf("preflight environment resolution failed: %w", err)
	}

	// Fetch raw service variables to inspect author templates
	rawVars := make(map[string]string)
	if deployer.store != nil {
		if appVars, err := deployer.store.ListServiceVariables(app.ID); err == nil {
			for _, v := range appVars {
				rawVars[v.Key] = v.Value
			}
		}
	}

	// Also check any service command / inline expressions
	if app.StartCommand != "" {
		rawVars["START_COMMAND"] = app.StartCommand
	}
	if app.BuildCommand != "" {
		rawVars["BUILD_COMMAND"] = app.BuildCommand
	}

	missing := FindMissingRequiredEnv(app.Name, rawVars, envVarsMap)
	if len(missing) > 0 {
		return FormatRequiredEnvError(missing)
	}

	// Test strict cross-service interpolation if applicable
	var registry map[string]map[string]string
	if deployer.ScopedEnvInterpolator != nil {
		registry, _ = deployer.ScopedEnvInterpolator(ctx, app.ProjectID, app.EnvironmentID)
	} else if deployer.EnvInterpolator != nil {
		registry, _ = deployer.EnvInterpolator(app.ProjectID)
	}

	if registry != nil {
		if _, err := build.InterpolateEnvVarsStrict(envVarsMap, registry); err != nil {
			return fmt.Errorf(`Service "%s" environment configuration error: %w`, app.Name, err)
		}
	}

	return nil
}
