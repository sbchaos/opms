package resource

type Schema []*Field

type Field struct {
	Name        string `mapstructure:"name,omitempty"`
	Type        string `mapstructure:"type,omitempty"`
	Description string `mapstructure:"description,omitempty"`

	// First label should be the primary label and others as extended
	Labels []string `mapstructure:"labels,omitempty"`

	Decimal *Decimal `mapstructure:"decimal,omitempty"`
}

type Decimal struct {
	Precision int32 `mapstructure:"precision"`
	Scale     int32 `mapstructure:"scale"`
}
