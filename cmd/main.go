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
		},
	}
	app.Flags = config.Flags
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
