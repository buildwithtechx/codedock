package projects

import (
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
	"context"
	"fmt"
	"strings"
)

func (sl *ServiceLinker) SetApplications(apps repositories.AppServiceRepository, variables repositories.ServiceVarRepository) {
	sl.apps, sl.variables = apps, variables
}

func (sl *ServiceLinker) GetNamespacedVariablesForEnvironment(ctx context.Context, project, environment string) (map[string]map[string]string, error) {
	registry := map[string]map[string]string{}
	databases, err := sl.databases.ListByProject(ctx, project)
	if err != nil {
		return nil, err
	}
	ambiguous := map[string]bool{}
	insert := func(name string, variables map[string]string) error {
		name = strings.ToLower(name)
		if ambiguous[name] {
			return nil
		}
		for existing := range registry {
			if strings.EqualFold(existing, name) {
				delete(registry, existing)
				ambiguous[name] = true
				return nil
			}
		}
		registry[name] = variables
		return nil
	}

	for _, database := range databases {
		if database.EnvironmentID == environment {
			registry[database.ID] = buildDatabaseEnvVars(database)
			if err := insert(database.Name, buildDatabaseEnvVars(database)); err != nil {
				return nil, err
			}
		}
	}
	if sl.apps != nil && sl.variables != nil {
		apps, err := sl.apps.ListByEnvironment(ctx, environment)
		if err != nil {
			return nil, err
		}
		for _, app := range apps {
			if app.ProjectID != project {
				continue
			}
			variables, err := sl.variables.ListByService(ctx, app.ID)
			if err != nil {
				return nil, err
			}
			values := map[string]string{"HOST": utils.NormalizeContainerName(app.ID), "PORT": fmt.Sprint(app.InternalPort)}
			if app.Replicas > 1 {
				values["HOST"] += "-1"
			}
			for _, variable := range variables {
				values[variable.Key] = variable.Value
			}
			registry[app.ID] = values
			if err := insert(app.Name, values); err != nil {
				return nil, err
			}
		}
	}
	return registry, nil
}

func (sl *ServiceLinker) GetLinkedVariablesForEnvironment(ctx context.Context, project, environment string) (map[string]string, error) {
	values := map[string]string{}
	databases, err := sl.databases.ListByProject(ctx, project)
	if err != nil {
		return nil, err
	}
	for _, database := range databases {
		if database.EnvironmentID != environment {
			continue
		}
		for key, value := range buildDatabaseEnvVars(database) {
			if _, duplicate := values[key]; !duplicate {
				values[key] = value
			}
		}
	}
	return values, nil
}
