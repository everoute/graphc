package informer

type Config struct {
	Resources    []string `yaml:"resources,omitempty"`
	Project      string   `yaml:"project,omitempty"`
	Package      string   `yaml:"package,omitempty"`
	SchemaModule string   `yaml:"schema_module,omitempty"`
	OutFile      string   `yaml:"outFile,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Package:      "informer",
		SchemaModule: "pkg/controller/tower/schema",
	}
}
