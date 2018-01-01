/*
Package config implements the functionality for processing bashly configuration files.
*/
package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/bryce/bashly/boxes"
)

// Config is the configuration format for the application.
type Config struct {
	LogsDirectory string         `json:"logsDirectory"`
	Boxes         []boxes.Config `json:"boxes"`
}

// Default loads the default application configuration.
func Default() *Config {
	cfg := &Config{}
	cfg.Boxes = []boxes.Config{}

	box := boxes.Config{}
	box.Name = "Script"
	box.Type = "Script"
	box.X0 = 0
	box.Y0 = 0
	box.X1 = 50
	box.Y1 = 100
	box.TabSize = 4
	cfg.Boxes = append(cfg.Boxes, box)

	box = boxes.Config{}
	box.Name = "Manual"
	box.Type = "Manual"
	box.RefName = "Script"
	box.X0 = 50
	box.Y0 = 0
	box.X1 = 100
	box.Y1 = 60
	cfg.Boxes = append(cfg.Boxes, box)

	box = boxes.Config{}
	box.Name = "Options"
	box.Type = "Options"
	box.RefName = "Script"
	box.X0 = 50
	box.Y0 = 60
	box.X1 = 100
	box.Y1 = 100
	cfg.Boxes = append(cfg.Boxes, box)

	return cfg
}

// Load loads the application configuration from a file.
func Load(cfgFile string) (*Config, error) {
	bytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(bytes, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
