package config

import "github.com/urfave/cli"

const (
	// uploader flags
	uploaderFlagDriver    = "driver"
	uploaderFlagRegion    = "region"
	uploaderFlagAccessKey = "access_key"
	uploaderFlagSecretKey = "secret_key"
	uploaderFlagBucket    = "bucket"
	uploaderFlagExclude   = "exclude"
	uploaderFlagSaveRoot  = "save_root"

	// uploader environments
	uploaderEnvAccessKey = "UPTOC_UPLOADER_AK"
	uploaderEnvSecretKey = "UPTOC_UPLOADER_SK"
)

// Flags defined the support flags for the cli
var Flags = []cli.Flag{
	cli.StringFlag{
		Name:  uploaderFlagDriver,
		Usage: "specify cloud storage engine",
		Value: "oss",
	},
	cli.StringFlag{
		Name:  uploaderFlagRegion,
		Usage: "specify region of the cloud platform",
	},
	cli.StringFlag{
		Name:  uploaderFlagBucket,
		Usage: "specify bucket name of the cloud platform",
	},
	cli.StringFlag{
		Name:   uploaderFlagAccessKey,
		Usage:  "specify key id of the cloud platform",
		EnvVar: uploaderEnvAccessKey,
	},
	cli.StringFlag{
		Name:   uploaderFlagSecretKey,
		Usage:  "specify key secret of the cloud platform",
		EnvVar: uploaderEnvSecretKey,
	},
	cli.StringFlag{
		Name:  uploaderFlagExclude,
		Usage: "specify exclude the given comma separated directories (example: --exclude=.cache,test)",
	},
	cli.StringFlag{
		Name:  uploaderFlagSaveRoot,
		Usage: "specify remote directory, default is root",
	},
}
