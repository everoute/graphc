package config

import (
	"bytes"
	"fmt"
	"os"

	gqlgenconfig "github.com/99designs/gqlgen/codegen/config"
	"github.com/everoute/graphc/tools/codegen/gqlgen"
	"github.com/everoute/graphc/tools/codegen/informer"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Project          string              `yaml:"project,omitempty"`
	SkipGqlGenerated bool                `yaml:"skipGqlGenerated,omitempty"`
	Informer         informer.Config     `yaml:"informer"`
	Gqlgen           gqlgenconfig.Config `yaml:"gqlgen"`
}

func GetConfig(configFile string) (*Config, error) {
	config := &Config{}
	config.SkipGqlGenerated = true
	config.Informer = *informer.DefaultConfig()
	gqlDefaultCfg, err := gqlgen.DefaultConfig()
	if err != nil {
		return nil, err
	}
	config.Gqlgen = *gqlDefaultCfg
	b, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)

	if err := dec.Decode(config); err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	if config.Informer.Project == "" {
		config.Informer.Project = config.Project
	}
	if config.Informer.SchemaModuleName == "" {
		config.Informer.SchemaModuleName = config.Gqlgen.Model.Package
	}

	err = gqlgen.CompleteConfig(&config.Gqlgen)

	return config, err
}

func (c *Config) GetInformerConfig() informer.Config {
	return c.Informer
}

func (c *Config) GetGqlgenConfig() gqlgenconfig.Config {
	return c.Gqlgen
}
