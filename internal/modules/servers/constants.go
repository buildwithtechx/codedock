package servers

type ServerStatus string

const (
	ServerStatusOnline       ServerStatus = "online"
	ServerStatusOffline      ServerStatus = "offline"
	ServerStatusProvisioning ServerStatus = "provisioning"
)

type ManagedProvider string

const (
	ManagedProviderHetzner ManagedProvider = "hetzner"
)

const (
	ManagedActionProvision = "provision"
	ManagedActionResize    = "resize"
	ManagedActionDelete    = "delete"
)

type TakeoverPlatform string

const (
	TakeoverPlatformDokploy TakeoverPlatform = "dokploy"
	TakeoverPlatformCoolify TakeoverPlatform = "coolify"
	TakeoverPlatformDokku   TakeoverPlatform = "dokku"
	TakeoverPlatformDocker  TakeoverPlatform = "docker"
)

type TakeoverStatus string

const (
	TakeoverStatusScanning TakeoverStatus = "scanning"
	TakeoverStatusScanned  TakeoverStatus = "scanned"
	TakeoverStatusAdopting TakeoverStatus = "adopting"
	TakeoverStatusDone     TakeoverStatus = "done"
	TakeoverStatusFailed   TakeoverStatus = "failed"
)
