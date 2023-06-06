package gqlgen

import (
	"os"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
)

const DefaultGeneratedFile = "generated.go"

func Run(cfg *config.Config, skipGqlGenerated bool) error {
	if err := api.Generate(cfg); err != nil {
		return err
	}

	if skipGqlGenerated {
		_, err := os.Stat(DefaultGeneratedFile)
		if err == nil {
			return os.Remove(DefaultGeneratedFile)
		}
	}
	return nil
}
