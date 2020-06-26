package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"

	"uptoc/config"
	"uptoc/engine"
	"uptoc/uploader"
)

const (
	// uploader flags
	uploaderFlagDriver    = "driver"
	uploaderFlagRegion    = "region"
	uploaderFlagAccessKey = "access_key"
	uploaderFlagSecretKey = "access_secret"
	uploaderFlagBucket    = "bucket"
	uploaderFlagExclude   = "exclude"

	// uploader environments
	uploaderEnvAccessKey = "UPTOC_UPLOADER_AK"
	uploaderEnvSecretKey = "UPTOC_UPLOADER_SK"
)

var (
	// RELEASE returns the release version
	release = "unknown"
	// REPO returns the git repository URL
	repo = "unknown"
	// COMMIT returns the short sha from git
	commit = "unknown"

	flags = []cli.Flag{
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
			Name:  uploaderFlagBucket,
			Usage: "specify bucket name of the cloud platform",
		},
		cli.StringFlag{
			Name:  uploaderFlagExclude,
			Usage: "specify exclude the given comma separated directories (example: --exclude=.cache,test)",
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "uptoc"
	app.Usage = "A cli tool to upload the dist file for the cloud engine."
	app.Copyright = "(c) 2019 saltbo.cn"
	app.Compiled = time.Now()
	app.Version = fmt.Sprintf("release: %s, repo: %s, commit: %s", release, repo, commit)
	app.Commands = cli.Commands{
		cli.Command{
			Name:   "config",
			Usage:  "config set up the engine for the bucket.",
			Action: configAction,
			Flags:  flags,
		},
	}
	app.Action = appAction
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func configAction(ctx *cli.Context) {
	c, err := config.Parse()
	if err != nil {
		log.Fatalln(err)
	}

	ak := ctx.String(uploaderFlagAccessKey)
	if ak != "" {
		c.Driver.AccessKey = ak
	}

	sk := ctx.String(uploaderFlagSecretKey)
	if sk != "" {
		c.Driver.SecretKey = sk
	}

	driver := ctx.String(uploaderFlagDriver)
	if driver != "" {
		c.Driver.Name = driver
	}

	region := ctx.String(uploaderFlagRegion)
	if region != "" {
		c.Driver.Region = region
	}

	bucket := ctx.String(uploaderFlagBucket)
	if bucket != "" {
		c.Driver.Bucket = bucket
	}

	exclude := ctx.String(uploaderFlagExclude)
	if exclude != "" {
		c.Core.Excludes = strings.Split(exclude, ",")
	}

	if err := c.Save(); err != nil {
		log.Fatalln(err)
	}
}

// 同步一个或多个文件或文件夹到指定目录，同步过程中会进行MD5判定，相同文件不再重复上传
func appAction(ctx *cli.Context) {
	conf, err := config.Parse()
	if err != nil {
		log.Fatalln(err)
	}

	ud, err := uploader.New(conf.Driver)
	if err != nil {
		log.Fatalln(err)
	}

	e := engine.New(conf.Core, ud)
	e.TailRun(ctx.Args()...)
}
