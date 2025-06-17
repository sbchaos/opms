package conf

const (
	DefaultFilename = "optimus.yaml"
)

type ClientConfig struct {
	Version    int          `yaml:"version"`
	Log        LogConfig    `yaml:"log"`
	Host       string       `yaml:"host"`
	Project    Project      `yaml:"project"`
	Namespaces []*Namespace `yaml:"namespaces"`
	Auth       Auth         `yaml:"auth"`
}

type LogConfig struct {
	Level  string `yaml:"level" default:"INFO"`
	Format string `yaml:"format"`
}

type Datastore struct {
	Type   string            `yaml:"type"`
	Path   string            `yaml:"path"`
	Backup map[string]string `yaml:"backup"`
}

type Project struct {
	Name        string            `yaml:"name"`
	Config      map[string]string `yaml:"config"`
	Variables   map[string]string `yaml:"variables"`
	PresetsPath string            `yaml:"preset_path"`
}

type Auth struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type Namespace struct {
	Name      string            `yaml:"name"`
	Config    map[string]string `yaml:"config"`
	Variables map[string]string `yaml:"variables"`
	Job       Job               `yaml:"job"`
	Datastore []Datastore       `yaml:"datastore"`
}

type Job struct {
	Path string `yaml:"path"`
}
