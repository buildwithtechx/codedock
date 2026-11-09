package apps

type AppServiceStatus string

const (
	AppServiceStatusCreated  AppServiceStatus = "created"
	AppServiceStatusBuilding AppServiceStatus = "building"
	AppServiceStatusStopped  AppServiceStatus = "stopped"
	AppServiceStatusRunning  AppServiceStatus = "running"
	AppServiceStatusError    AppServiceStatus = "error"
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
	RuntimeModeStandard   RuntimeMode = "standard"
	RuntimeModeServerless RuntimeMode = "serverless"
	RuntimeModeStatic     RuntimeMode = "static"
	RuntimeModeWorker     RuntimeMode = "worker"
	RuntimeModeWeb        RuntimeMode = "web"
)

type DrainType string

const (
	DrainTypeHTTP   DrainType = "http"
	DrainTypeSyslog DrainType = "syslog"
)
