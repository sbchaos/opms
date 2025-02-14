package resource

type ExternalTable struct {
	Name string

	Description string          `mapstructure:"description,omitempty"`
	Schema      Schema          `mapstructure:"schema,omitempty"`
	Source      *ExternalSource `mapstructure:"source,omitempty"`

	Hints       map[string]string      `mapstructure:"hints,omitempty"`
	ExtraConfig map[string]interface{} `mapstructure:",remain"`
}

func (e *ExternalTable) FullName() string {
	return e.Name
}

type ExternalSource struct {
	SourceType string   `mapstructure:"type,omitempty"`
	SourceURIs []string `mapstructure:"uris,omitempty"`

	// Additional configs for CSV, GoogleSheets, LarkSheets formats.
	Config map[string]interface{} `mapstructure:"config"`

	SerdeProperties map[string]string `mapstructure:"serde_properties"`
	TableProperties map[string]string `mapstructure:"table_properties"`

	SyncInterval int64    `mapstructure:"sync_interval_in_hrs,omitempty"`
	Jars         []string `mapstructure:"jars,omitempty"`
	Location     string   `mapstructure:"location,omitempty"`
	Range        string   `mapstructure:"range,omitempty"`
}
