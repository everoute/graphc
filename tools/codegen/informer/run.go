package informer

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed factory.tpl
var factoryTmpl string

func Run(cfg *Config) error {
	tmpl, err := template.New("factory").Parse(factoryTmpl)
	if err != nil {
		return err
	}

	_, err = os.Stat(cfg.OutFile)
	if err == nil {
		if err := os.Remove(cfg.OutFile); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	file, err := os.OpenFile(cfg.OutFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	err = tmpl.Execute(file, cfg)
	if err != nil {
		return err
	}

	return nil
}
