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

// Config provides yml configuration.
type Config struct {
	f *os.File

	Core   engine.Config   `yaml:"core"`
	Driver uploader.Config `yaml:"driver"`
}

// Parse returns Config from RC file if no flags and
// returns Config from flags if any flags exist.
func Parse(ctx *cli.Context) (*Config, error) {
	if ctx.NumFlags() > 0 {
		c := &Config{
			Core: engine.Config{
				SaveRoot:  ctx.String(uploaderFlagSaveRoot),
				ForceSync: true,
			},
			Driver: uploader.Config{
				Name:      ctx.String(uploaderFlagDriver),
				Region:    ctx.String(uploaderFlagRegion),
				Bucket:    ctx.String(uploaderFlagBucket),
				AccessKey: ctx.String(uploaderFlagAccessKey),
				SecretKey: ctx.String(uploaderFlagSecretKey),
			},
		}
		exclude := ctx.String(uploaderFlagExclude)
		if exclude != "" {
			c.Core.Excludes = strings.Split(exclude, ",")
		}

		return c, nil
	}

	return ParseFromRC()
}

// ParseFromRC returns Config from rc file
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

// Prompt implement a prompt for the config
func (c *Config) Prompt() error {
	prompts := []struct {
		label    string
		value    *string
		mask     rune
		validate promptui.ValidateFunc
	}{
		{label: uploaderFlagDriver, value: &c.Driver.Name, validate: uploader.DriverValidate},
		{label: uploaderFlagRegion, value: &c.Driver.Region},
		{label: uploaderFlagBucket, value: &c.Driver.Bucket},
		{label: uploaderFlagAccessKey, value: &c.Driver.AccessKey},
		{label: uploaderFlagSecretKey, value: &c.Driver.SecretKey, mask: '*'},
		{label: uploaderFlagSaveRoot, value: &c.Core.SaveRoot},
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

	defer c.f.Close()
	c.f.Seek(0, io.SeekStart)
	return yaml.NewEncoder(c.f).Encode(c)
}
