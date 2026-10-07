package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"fmt"
	"path"
)

func postgresObjects(plan *models.ClusterDataPlan, metadata func(string) map[string]any, destination, source *models.S3Destination, sourceName string) (map[string]any, []map[string]any, error) {
	spec := plan.Record.Spec
	_, name := DataIdentity(spec)
	storage := map[string]any{"size": fmt.Sprintf("%dGi", spec.StorageGiB)}
	if spec.StorageClass != "" {
		storage["storageClass"] = spec.StorageClass
	}
	settings := map[string]any{"instances": spec.Instances, "imageName": spec.Image, "storage": storage, "enableSuperuserAccess": false, "bootstrap": map[string]any{"initdb": map[string]any{"database": "app", "owner": "app", "secret": map[string]string{"name": name + "-auth"}}}, "affinity": map[string]any{"enablePodAntiAffinity": true, "podAntiAffinityType": "required", "topologyKey": "kubernetes.io/hostname"}}
	settings["inheritedMetadata"] = map[string]any{"labels": metadata(name)["labels"]}
	if spec.Synchronous {
		settings["postgresql"] = map[string]any{"parameters": map[string]string{"synchronous_commit": "remote_apply", "synchronous_standby_names": "ANY 1 (*)"}}
	}
	extra := []map[string]any{}
	objectStore := func(suffix string, dest *models.S3Destination, serverName string) {
		credentials := name + suffix + "-credentials"
		storeName := name + suffix
		extra = append(extra, map[string]any{"apiVersion": "v1", "kind": "Secret", "metadata": metadata(credentials), "type": "Opaque", "stringData": map[string]string{"ACCESS_KEY_ID": dest.AccessKeyID, "ACCESS_SECRET_KEY": dest.SecretAccessKey}})
		configuration := map[string]any{"destinationPath": "s3://" + dest.Bucket + "/" + path.Join(dest.PathPrefix, "codedock", plan.Record.ProjectID, serverName), "s3Credentials": map[string]any{"accessKeyId": map[string]string{"name": credentials, "key": "ACCESS_KEY_ID"}, "secretAccessKey": map[string]string{"name": credentials, "key": "ACCESS_SECRET_KEY"}}, "wal": map[string]string{"compression": "gzip"}, "data": map[string]string{"compression": "gzip"}}
		if dest.Endpoint != "" {
			configuration["endpointURL"] = dest.Endpoint
		}
		extra = append(extra, map[string]any{"apiVersion": "barmancloud.cnpg.io/v1", "kind": "ObjectStore", "metadata": metadata(storeName), "spec": map[string]any{"configuration": configuration, "retentionPolicy": "30d"}})
	}
	if destination != nil {
		objectStore("-backups", destination, name)
		settings["plugins"] = []map[string]any{{"name": "barman-cloud.cloudnative-pg.io", "isWALArchiver": true, "parameters": map[string]string{"barmanObjectName": name + "-backups"}}}
		extra = append(extra, map[string]any{"apiVersion": "postgresql.cnpg.io/v1", "kind": "ScheduledBackup", "metadata": metadata(name + "-daily"), "spec": map[string]any{"schedule": "0 0 2 * * *", "backupOwnerReference": "cluster", "cluster": map[string]string{"name": name}, "method": "plugin", "pluginConfiguration": map[string]string{"name": "barman-cloud.cloudnative-pg.io"}}})
	}
	if source != nil {
		if sourceName == "" || plan.RestoreTime == "" {
			return nil, nil, fmt.Errorf("restore requires a source and reviewed recovery time")
		}
		objectStore("-source", source, sourceName)
		settings["bootstrap"] = map[string]any{"recovery": map[string]any{"source": "source", "database": "app", "owner": "app", "secret": map[string]string{"name": name + "-auth"}, "recoveryTarget": map[string]string{"targetTime": plan.RestoreTime}}}
		settings["externalClusters"] = []map[string]any{{"name": "source", "plugin": map[string]any{"name": "barman-cloud.cloudnative-pg.io", "parameters": map[string]string{"barmanObjectName": name + "-source", "serverName": sourceName}}}}
	}
	return map[string]any{"apiVersion": "postgresql.cnpg.io/v1", "kind": "Cluster", "metadata": metadata(name), "spec": settings}, extra, nil
}
