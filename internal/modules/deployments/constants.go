package deployments

type DeploymentStatus string

const (
	DeploymentStatusPending        DeploymentStatus = "pending"
	DeploymentStatusPreparing      DeploymentStatus = "PREPARING"
	DeploymentStatusCloning        DeploymentStatus = "CLONING"
	DeploymentStatusPulling        DeploymentStatus = "PULLING"
	DeploymentStatusBuilding       DeploymentStatus = "BUILDING"
	DeploymentStatusStarting       DeploymentStatus = "STARTING"
	DeploymentStatusReadiness      DeploymentStatus = "READINESS"
	DeploymentStatusRouting        DeploymentStatus = "ROUTING"
	DeploymentStatusReady          DeploymentStatus = "READY"
	DeploymentStatusActive         DeploymentStatus = "ACTIVE"
	DeploymentStatusSuccess        DeploymentStatus = "SUCCESS"
	DeploymentStatusFailed         DeploymentStatus = "FAILED"
	DeploymentStatusCancelled      DeploymentStatus = "CANCELLED"
	DeploymentStatusRemoved        DeploymentStatus = "REMOVED"
	DeploymentStatusSlept          DeploymentStatus = "SLEPT"
	DeploymentStatusPartialFailure DeploymentStatus = "PARTIAL_FAILURE"
)

type BuildEngine string

const (
	BuildEngineAuto       BuildEngine = "auto"
	BuildEngineDockerfile BuildEngine = "dockerfile"
	BuildEngineNixpacks   BuildEngine = "nixpacks"
	BuildEngineBuildpacks BuildEngine = "buildpacks"
	BuildEngineRailpack   BuildEngine = "railpack"
	BuildEngineServerless BuildEngine = "serverless"
)

type RuntimeMode string

const (
	RuntimeModeWeb    RuntimeMode = "web"
	RuntimeModeWorker RuntimeMode = "worker"
	RuntimeModeStatic RuntimeMode = "static"
	RuntimeModeDocker RuntimeMode = "docker"
)

const (
	DefaultMaxConcurrentBuilds = 3
	DefaultServiceBuildLimit   = 1
	DefaultRetentionDays       = 30
	DefaultRetentionCount      = 10
	MaxArchiveUploadSizeBytes  = 100 * 1024 * 1024 // 100 MB
)
