package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"uptoc/engine"
	"uptoc/uploader"
)

type Config struct {
	f *os.File

	Core   engine.Config   `yaml:"core"`
	Driver uploader.Config `yaml:"driver"`
}

func newConfig(f *os.File) *Config {
	return &Config{f: f}
}

func Parse() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	rcPath := filepath.Join(homeDir, ".uptocrc")
	f, err := os.OpenFile(rcPath, os.O_CREATE|os.O_RDWR, 0644)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("please setup your config by run `uptoc config`")
	} else if err != nil {
		return nil, fmt.Errorf("open .uptocrc failed: %s", err)
	}

	c := newConfig(f)
	yd := yaml.NewDecoder(f)
	if err := yd.Decode(c); err != nil {
		return nil, err
	}

	if strings.HasPrefix(c.Core.SaveRoot, "/") {
		c.Core.SaveRoot = c.Core.SaveRoot[1:]
	}

	if !strings.HasSuffix(c.Core.VisitHost, "/") {
		c.Core.VisitHost += "/"
	}

	return c, nil
}

func (c *Config) Save() error {
	c.f.Seek(0, 0)
	ye := yaml.NewEncoder(c.f)
	if err := ye.Encode(c); err != nil {
		return err
	}

	return nil
}
