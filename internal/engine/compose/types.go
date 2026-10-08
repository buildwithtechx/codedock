package compose

type ComposeTemplate struct {
	Version   string                    `yaml:"version,omitempty"`
	Services  map[string]ComposeService `yaml:"services"`
	Volumes   map[string]any            `yaml:"volumes,omitempty"`
	XCodedock *CodedockMetadata         `yaml:"x-codedock,omitempty"`
}

type ComposeService struct {
	Image       string              `yaml:"image"`
	Environment []string            `yaml:"environment,omitempty"`
	Ports       []string            `yaml:"ports,omitempty"`
	Volumes     []string            `yaml:"volumes,omitempty"`
	Command     []string            `yaml:"command,omitempty"`
	DependsOn   []string            `yaml:"depends_on,omitempty"`
	Healthcheck *CatalogHealthcheck `yaml:"healthcheck,omitempty"`
	XCodedock   *CodedockMetadata   `yaml:"x-codedock,omitempty"`
}

type CatalogHealthcheck struct {
	Test        []string `yaml:"test"`
	Interval    string   `yaml:"interval,omitempty"`
	Timeout     string   `yaml:"timeout,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod string   `yaml:"start_period,omitempty"`
}

type CodedockMetadata struct {
	IsDatabase       bool                     `yaml:"is_database,omitempty"`
	IsOneClick       bool                     `yaml:"is_one_click,omitempty"`
	Verified         bool                     `yaml:"verified,omitempty"`
	Name             string                   `yaml:"name,omitempty"`
	Description      string                   `yaml:"description,omitempty"`
	Category         string                   `yaml:"category,omitempty"`
	Icon             string                   `yaml:"icon,omitempty"`
	ConnectionString string                   `yaml:"connection_string,omitempty"`
	Secrets          []CodedockSecretSpec     `yaml:"secrets,omitempty"`
	Inputs           []CodedockInputSpec      `yaml:"inputs,omitempty"`
	Backup           *CodedockBackupMetadata  `yaml:"backup,omitempty"`
	Restore          *CodedockRestoreMetadata `yaml:"restore,omitempty"`
}

type CodedockSecretSpec struct {
	Var    string `yaml:"var"`
	Label  string `yaml:"label,omitempty"`
	Length int    `yaml:"length,omitempty"`
}

type CodedockInputSpec struct {
	Var      string `yaml:"var"`
	Label    string `yaml:"label,omitempty"`
	Default  string `yaml:"default,omitempty"`
	Required bool   `yaml:"required,omitempty"`
}

type CodedockBackupMetadata struct {
	Command       []string `yaml:"command,omitempty"`
	FileExtension string   `yaml:"file_extension,omitempty"`
}

type CodedockRestoreMetadata struct {
	Command []string `yaml:"command,omitempty"`
}
