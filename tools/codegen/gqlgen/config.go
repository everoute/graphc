package gqlgen

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/99designs/gqlgen/codegen/config"
	"gopkg.in/yaml.v3"
)

//go:embed gqlgen.yaml
var gqlConfigFile embed.FS

func DefaultConfig() (*config.Config, error) {
	resConfig := config.DefaultConfig()
	b, err := gqlConfigFile.ReadFile("gqlgen.yaml")
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)

	if err := dec.Decode(resConfig); err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	if err := config.CompleteConfig(resConfig); err != nil {
		return nil, err
	}
	return resConfig, nil
}
