package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"uptoc/engine"
	"uptoc/uploader"
)

const (
	// uploader flags
	UploaderFlagDriver    = "driver"
	UploaderFlagRegion    = "region"
	UploaderFlagAccessKey = "access_key"
	UploaderFlagSecretKey = "secret_key"
	UploaderFlagBucket    = "bucket"
	UploaderFlagExclude   = "exclude"
	UploaderFlagSaveRoot  = "save_root"

	// uploader environments
	UploaderEnvAccessKey = "UPTOC_UPLOADER_AK"
	UploaderEnvSecretKey = "UPTOC_UPLOADER_SK"
)

type Config struct {
	f *os.File

	Core   engine.Config   `yaml:"core"`
	Driver uploader.Config `yaml:"driver"`
}

func Parse(ctx *cli.Context) (*Config, error) {
	if ctx.NumFlags() > 0 {
		c := &Config{
			Core: engine.Config{
				SaveRoot:  ctx.String(UploaderFlagSaveRoot),
				ForceSync: true,
			},
			Driver: uploader.Config{
				Name:      ctx.String(UploaderFlagDriver),
				Region:    ctx.String(UploaderFlagRegion),
				Bucket:    ctx.String(UploaderFlagBucket),
				AccessKey: ctx.String(UploaderFlagAccessKey),
				SecretKey: ctx.String(UploaderFlagSecretKey),
			},
		}
		exclude := ctx.String(UploaderFlagExclude)
		if exclude != "" {
			c.Core.Excludes = strings.Split(exclude, ",")
		}

		return c, nil
	}

	return ParseFromRC()
}

func ParseFromRC() (*Config, error) {
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

	c := &Config{f: f}
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
	c.f.Seek(0, io.SeekStart)
	defer c.f.Close()

	ye := yaml.NewEncoder(c.f)
	if err := ye.Encode(c); err != nil {
		return err
	}

	return nil
}

func (c *Config) Prompt() error {
	prompts := []struct {
		label    string
		value    *string
		mask     rune
		validate promptui.ValidateFunc
	}{
		{label: UploaderFlagDriver, value: &c.Driver.Name, validate: uploader.DriverValidate},
		{label: UploaderFlagRegion, value: &c.Driver.Region},
		{label: UploaderFlagBucket, value: &c.Driver.Bucket},
		{label: UploaderFlagAccessKey, value: &c.Driver.AccessKey},
		{label: UploaderFlagSecretKey, value: &c.Driver.SecretKey, mask: '*'},
		{label: UploaderFlagSaveRoot, value: &c.Core.SaveRoot},
	}

	for _, prompt := range prompts {
		pp := promptui.Prompt{
			Label:    prompt.label,
			Default:  *prompt.value,
			Validate: prompt.validate,
			Mask:     prompt.mask,
		}

		value, err := pp.Run()
		if err != nil {
			return err
		}

		*prompt.value = value
	}

	return nil
}
