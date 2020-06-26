package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	"uptoc/config"
	"uptoc/engine"
	"uptoc/uploader"
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
			Name:  config.UploaderFlagDriver,
			Usage: "specify cloud storage engine",
			Value: "oss",
		},
		cli.StringFlag{
			Name:  config.UploaderFlagRegion,
			Usage: "specify region of the cloud platform",
		},
		cli.StringFlag{
			Name:   config.UploaderFlagAccessKey,
			Usage:  "specify key id of the cloud platform",
			EnvVar: config.UploaderEnvAccessKey,
		},
		cli.StringFlag{
			Name:   config.UploaderFlagSecretKey,
			Usage:  "specify key secret of the cloud platform",
			EnvVar: config.UploaderEnvSecretKey,
		},
		cli.StringFlag{
			Name:  config.UploaderFlagBucket,
			Usage: "specify bucket name of the cloud platform",
		},
		cli.StringFlag{
			Name:  config.UploaderFlagExclude,
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
	app.Flags = flags
	app.Action = appAction
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func configAction(ctx *cli.Context) {
	c, err := config.ParseFromRC()
	if err != nil {
		log.Fatalln(err)
	}

	if err := c.Prompt(); err != nil {
		log.Fatalln(err)
	}

	if err := c.Save(); err != nil {
		log.Fatalln(err)
	}
}

func appAction(ctx *cli.Context) {
	conf, err := config.Parse(ctx)
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
