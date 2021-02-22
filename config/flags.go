package config

import "github.com/urfave/cli"

const (
	// uploader flags
	uploaderDriver        = "driver"
	uploaderRegion        = "region"
	uploaderAccessKey     = "access_key"
	uploaderSecretKey     = "secret_key"
	uploaderBucket        = "bucket"
	uploaderExclude       = "exclude"
	uploaderSaveRoot      = "save_root"
	uploaderVisitHost     = "visit_host"
	uploaderRemoteExclude = "remote_exclude"

	// uploader environments
	uploaderEnvAccessKey = "UPTOC_UPLOADER_AK"
	uploaderEnvSecretKey = "UPTOC_UPLOADER_SK"
)

// Flags defined the support flags for the cli
var Flags = []cli.Flag{
	cli.StringFlag{
		Name:  uploaderDriver,
		Usage: "specify cloud storage engine",
		Value: "oss",
	},
	cli.StringFlag{
		Name:  uploaderRegion,
		Usage: "specify region of the cloud platform",
	},
	cli.StringFlag{
		Name:  uploaderBucket,
		Usage: "specify bucket name of the cloud platform",
	},
	cli.StringFlag{
		Name:   uploaderAccessKey,
		Usage:  "specify key id of the cloud platform",
		EnvVar: uploaderEnvAccessKey,
	},
	cli.StringFlag{
		Name:   uploaderSecretKey,
		Usage:  "specify key secret of the cloud platform",
		EnvVar: uploaderEnvSecretKey,
	},
	cli.StringFlag{
		Name:  uploaderExclude,
		Usage: "specify exclude the given comma separated directories (example: --exclude=.cache,test)",
	},
	cli.StringFlag{
		Name:  uploaderSaveRoot,
		Usage: "specify remote directory, default is root",
	},
	cli.StringFlag{
		Name:  uploaderRemoteExclude,
		Usage: "specify exclude the given comma separated remote directories (example: --exclude=.cache,test)",
	},
}
