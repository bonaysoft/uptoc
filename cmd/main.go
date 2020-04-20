package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"

	"uptoc/core"
	"uptoc/uploader"
)

const (
	// uploader configs from os envs
	EnvUploaderAccessKey    = "UPLOADER_ACCESS_KEY"
	EnvUploaderAccessSecret = "UPLOADER_ACCESS_SECRET"

	// config from cmd flags
	UploaderDriver       = "driver"
	UploaderEndpoint     = "endpoint"
	UploaderAccessKey    = "access_key"
	UploaderAccessSecret = "access_secret"
	UploaderBucket       = "bucket"
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
			Name:  UploaderDriver,
			Usage: "specify cloud storage engine",
			Value: "oss",
		},
		cli.StringFlag{
			Name:     UploaderEndpoint,
			Usage:    "specify endpoint of the cloud platform",
			Required: true,
		},
		cli.StringFlag{
			Name:     UploaderAccessKey,
			Usage:    "specify key id of the cloud platform",
			EnvVar:   EnvUploaderAccessKey,
			Required: true,
		},
		cli.StringFlag{
			Name:     UploaderAccessSecret,
			Usage:    "specify key secret of the cloud platform",
			EnvVar:   EnvUploaderAccessSecret,
			Required: true,
		},
		cli.StringFlag{
			Name:     UploaderBucket,
			Usage:    "specify bucket name of the cloud platform",
			Required: true,
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
	app.Flags = flags
	app.Action = action
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) {
	driver := c.String(UploaderDriver)
	endpoint := c.String(UploaderEndpoint)
	accessKey := c.String(UploaderAccessKey)
	accessSecret := c.String(UploaderAccessSecret)
	bucketName := c.String(UploaderBucket)
	uploadDriver, err := uploader.New(driver, endpoint, accessKey, accessSecret, bucketName)
	if err != nil {
		log.Fatalln(err)
	}

	dirPath := c.Args().First()
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}

	e := core.NewEngine(uploadDriver)
	if err := e.LoadAndCompareObjects(dirPath); err != nil {
		log.Fatalln(err)
	}

	if err := e.Sync(); err != nil {
		log.Fatalln(err)
	}
}
