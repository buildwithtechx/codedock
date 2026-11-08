package clusters

type ClusterStatus string

const (
	ClusterStatusReady        ClusterStatus = "ready"
	ClusterStatusProvisioning ClusterStatus = "provisioning"
	ClusterStatusFailed       ClusterStatus = "failed"
	ClusterStatusDeleting     ClusterStatus = "deleting"
)

type ClusterDataStatus string

const (
	ClusterDataStatusReady       ClusterDataStatus = "ready"
	ClusterDataStatusApplying    ClusterDataStatus = "applying"
	ClusterDataStatusFailed      ClusterDataStatus = "failed"
	ClusterDataStatusTerminating ClusterDataStatus = "terminating"
)
