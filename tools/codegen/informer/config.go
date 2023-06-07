package informer

type Config struct {
	Resources        []string `yaml:"resources,omitempty"`
	Project          string   `yaml:"project,omitempty"`
	Package          string   `yaml:"package,omitempty"`
	SchemaModulePath string   `yaml:"schemaModulePath,omitempty"`
	SchemaModuleName string   `yaml:"schemaModuleName,omitempty"`
	OutFile          string   `yaml:"outFile,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Package: "informer",
	}
}
