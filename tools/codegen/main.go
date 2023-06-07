package main

import (
	"flag"
	"fmt"

	"github.com/everoute/graphc/tools/codegen/config"
	"github.com/everoute/graphc/tools/codegen/gqlgen"
	"github.com/everoute/graphc/tools/codegen/informer"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "codegen config file")
	flag.Parse()

	if configFile == "" {
		flag.Usage()
	}
	cfg, err := config.GetConfig(configFile)
	if err != nil {
		panic(fmt.Errorf("failed to get config: %v", err))
	}

	gqlCfg := cfg.GetGqlgenConfig()
	err = gqlgen.Run(&gqlCfg, cfg.SkipGqlGenerated)
	if err != nil {
		panic(fmt.Errorf("failed to gen go types code: %v", err))
	}

	infCfg := cfg.GetInformerConfig()
	err = informer.Run(&infCfg)
	if err != nil {
		panic(fmt.Errorf("failed to gen informer code: %v", err))
	}
}
