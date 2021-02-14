package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"uptoc/engine"
	"uptoc/uploader"
)

// Config provides yml configuration.
type Config struct {
	Core   engine.Config   `yaml:"core"`
	Driver uploader.Config `yaml:"driver"`
}

func New() *Config {
	return &Config{}
}

// Parse returns Config from RC file if no flags and
// returns Config from flags if any flags exist.
func NewWithCtx(ctx *cli.Context) (*Config, error) {
	c := New()
	if ctx.NumFlags() > 0 {
		c.Core = engine.Config{
			SaveRoot:  ctx.String(uploaderSaveRoot),
			ForceSync: true,
		}
		c.Driver = uploader.Config{
			Name:      ctx.String(uploaderDriver),
			Region:    ctx.String(uploaderRegion),
			Bucket:    ctx.String(uploaderBucket),
			AccessKey: ctx.String(uploaderAccessKey),
			SecretKey: ctx.String(uploaderSecretKey),
		}
		if exclude := ctx.String(uploaderExclude); exclude != "" {
			c.Core.Excludes = strings.Split(exclude, ",")
		}
		if ossExclude := ctx.String(uploaderOssExclude); ossExclude != "" {
			c.Core.OssExcludes = strings.Split(ossExclude, ",")
		}
	} else if err := c.Parse(); err != nil {
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

// Parse returns Config from rc file
func (c *Config) Parse() error {
	f, err := c.open(os.O_RDONLY)
	if err != nil {
		return err
	}

	return yaml.NewDecoder(f).Decode(c)
}

// Prompt implement a prompt for the config
func (c *Config) Prompt() error {
	prompts := []struct {
		label    string
		value    *string
		mask     rune
		validate promptui.ValidateFunc
	}{
		{label: uploaderDriver, value: &c.Driver.Name, validate: uploader.DriverValidate},
		{label: uploaderRegion, value: &c.Driver.Region},
		{label: uploaderBucket, value: &c.Driver.Bucket},
		{label: uploaderAccessKey, value: &c.Driver.AccessKey},
		{label: uploaderSecretKey, value: &c.Driver.SecretKey, mask: '*'},
		{label: uploaderSaveRoot, value: &c.Core.SaveRoot},
		{label: uploaderVisitHost, value: &c.Core.VisitHost},
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

	f, err := c.open(os.O_CREATE | os.O_WRONLY)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(f).Encode(c)
}

func (c *Config) open(flag int) (*os.File, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	rcPath := filepath.Join(homeDir, ".uptocrc")
	f, err := os.OpenFile(rcPath, flag, 0644)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("please setup your config by run `uptoc config`")
	} else if err != nil {
		return nil, fmt.Errorf("open .uptocrc failed: %s", err)
	}

	return f, err
}
