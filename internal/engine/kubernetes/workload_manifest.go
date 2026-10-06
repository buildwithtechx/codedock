package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func WorkloadManifest(workload *models.KubernetesWorkload) (string, error) {
	if err := ValidateWorkload(workload); err != nil {
		return "", err
	}
	app, target := &workload.App, &workload.Target
	namespace, name, err := WorkloadIdentity(app)
	if err != nil {
		return "", err
	}
	labels := map[string]string{"codedock.run/service": app.ID, "codedock.run/project": app.ProjectID, "codedock.run/environment": app.EnvironmentID}
	metadata := func(resource string) map[string]any {
		return map[string]any{"name": resource, "namespace": namespace, "labels": labels}
	}
	object := func(api, kind, resource string, spec map[string]any) map[string]any {
		return map[string]any{"apiVersion": api, "kind": kind, "metadata": metadata(resource), "spec": spec}
	}
	objects := []map[string]any{{"apiVersion": "v1", "kind": "Namespace", "metadata": map[string]any{"name": namespace, "labels": map[string]string{"codedock.run/project": app.ProjectID, "codedock.run/environment": app.EnvironmentID}}}}
	objects = append(objects, map[string]any{"apiVersion": "v1", "kind": "Secret", "metadata": metadata(name + "-env"), "type": "Opaque", "stringData": workload.Variables})
	container := map[string]any{"name": "app", "image": workload.Image, "imagePullPolicy": "Always", "ports": []map[string]any{{"containerPort": app.InternalPort}}, "envFrom": []map[string]any{{"secretRef": map[string]string{"name": name + "-env"}}}}
	if app.StartCommand != "" {
		container["command"] = []string{"sh", "-lc", app.StartCommand}
	}
	resources := map[string]any{}
	limits := map[string]string{}
	if app.CPULimit > 0 {
		limits["cpu"] = fmt.Sprintf("%dm", int(app.CPULimit*1000))
	}
	if app.MemoryLimit > 0 {
		limits["memory"] = fmt.Sprintf("%dMi", app.MemoryLimit)
	}
	if len(limits) > 0 {
		resources["limits"] = limits
		resources["requests"] = limits
		container["resources"] = resources
	}
	probe := map[string]any{"initialDelaySeconds": 2, "periodSeconds": 5, "timeoutSeconds": 3, "failureThreshold": 12}
	if app.HealthCheckPath != "" {
		probe["httpGet"] = map[string]any{"path": app.HealthCheckPath, "port": app.InternalPort}
	} else {
		probe["tcpSocket"] = map[string]any{"port": app.InternalPort}
	}
	if app.RuntimeMode != models.RuntimeModeWorker {
		container["readinessProbe"] = probe
	} else {
		delete(container, "ports")
	}
	pod := map[string]any{"containers": []map[string]any{container}, "terminationGracePeriodSeconds": 30}
	if len(target.NodeIDs) > 0 {
		nodes := append([]string{}, target.NodeIDs...)
		sort.Strings(nodes)
		for i := range nodes {
			nodes[i] = "codedock-" + nodes[i]
		}
		pod["affinity"] = map[string]any{"nodeAffinity": map[string]any{"requiredDuringSchedulingIgnoredDuringExecution": map[string]any{"nodeSelectorTerms": []map[string]any{{"matchExpressions": []map[string]any{{"key": "kubernetes.io/hostname", "operator": "In", "values": nodes}}}}}}}
	}
	mounts, volumes := []map[string]any{}, []map[string]any{}
	for _, volume := range target.Volumes {
		claim := name + "-" + volume.Name
		mode := "ReadWriteOnce"
		if volume.Shared {
			mode = "ReadWriteMany"
		}
		spec := map[string]any{"accessModes": []string{mode}, "resources": map[string]any{"requests": map[string]string{"storage": fmt.Sprintf("%dGi", volume.SizeGiB)}}}
		if volume.StorageClass != "" {
			spec["storageClassName"] = volume.StorageClass
		}
		objects = append(objects, object("v1", "PersistentVolumeClaim", claim, spec))
		volumes = append(volumes, map[string]any{"name": volume.Name, "persistentVolumeClaim": map[string]string{"claimName": claim}})
		mounts = append(mounts, map[string]any{"name": volume.Name, "mountPath": volume.MountPath})
	}
	if len(volumes) > 0 {
		pod["volumes"] = volumes
		container["volumeMounts"] = mounts
	}
	if workload.Registry != nil {
		registry := workload.Registry
		host := strings.TrimSuffix(registry.RegistryURL, "/")
		if strings.Contains(host, "://") {
			parsed, err := url.Parse(host)
			if err != nil {
				return "", err
			}
			host = parsed.Host
		}
		authentication, err := json.Marshal(map[string]any{"auths": map[string]any{host: map[string]string{"username": registry.Username, "password": registry.PasswordToken}}})
		if err != nil {
			return "", err
		}
		objects = append(objects, map[string]any{"apiVersion": "v1", "kind": "Secret", "metadata": metadata(name + "-registry"), "type": "kubernetes.io/dockerconfigjson", "stringData": map[string]string{".dockerconfigjson": string(authentication)}})
		pod["imagePullSecrets"] = []map[string]string{{"name": name + "-registry"}}
	}
	strategy := map[string]any{"type": "RollingUpdate", "rollingUpdate": map[string]any{"maxSurge": 1, "maxUnavailable": 0}}
	if len(target.Volumes) > 0 {
		strategy = map[string]any{"type": "Recreate"}
	}
	objects = append(objects, object("apps/v1", "Deployment", name, map[string]any{"replicas": app.Replicas, "revisionHistoryLimit": 5, "progressDeadlineSeconds": 300, "strategy": strategy, "selector": map[string]any{"matchLabels": map[string]string{"codedock.run/service": app.ID}}, "template": map[string]any{"metadata": map[string]any{"labels": labels}, "spec": pod}}))
	if app.RuntimeMode != models.RuntimeModeWorker {
		objects = append(objects, object("v1", "Service", name, map[string]any{"selector": map[string]string{"codedock.run/service": app.ID}, "ports": []map[string]any{{"port": 80, "targetPort": app.InternalPort}}}))
	}
	if app.Domain != "" && app.RuntimeMode == models.RuntimeModeWeb {
		host := app.Domain
		if strings.Contains(host, "://") {
			parsed, err := url.Parse(host)
			if err != nil {
				return "", err
			}
			host = parsed.Hostname()
		}
		objects = append(objects, object("networking.k8s.io/v1", "Ingress", name, map[string]any{"ingressClassName": "traefik", "rules": []map[string]any{{"host": host, "http": map[string]any{"paths": []map[string]any{{"path": "/", "pathType": "Prefix", "backend": map[string]any{"service": map[string]any{"name": name, "port": map[string]int{"number": 80}}}}}}}}}))
	}
	data, err := json.Marshal(map[string]any{"apiVersion": "v1", "kind": "List", "items": objects})
	if err != nil {
		return "", err
	}
	return string(data), nil
}
