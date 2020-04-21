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
	// uploader flags
	uploaderFlagDriver       = "driver"
	uploaderFlagEndpoint     = "endpoint"
	uploaderFlagAccessKey    = "access_key"
	uploaderFlagAccessSecret = "access_secret"
	uploaderFlagBucket       = "bucket"

	// uploader environments
	uploaderEnvAccessKey    = "uploaderFlag_ACCESS_KEY"
	uploaderEnvAccessSecret = "uploaderFlag_ACCESS_SECRET"
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
			Name:     uploaderFlagEndpoint,
			Usage:    "specify endpoint of the cloud platform",
			Required: true,
		},
		cli.StringFlag{
			Name:     uploaderFlagAccessKey,
			Usage:    "specify key id of the cloud platform",
			EnvVar:   uploaderEnvAccessKey,
			Required: true,
		},
		cli.StringFlag{
			Name:     uploaderFlagAccessSecret,
			Usage:    "specify key secret of the cloud platform",
			EnvVar:   uploaderEnvAccessSecret,
			Required: true,
		},
		cli.StringFlag{
			Name:     uploaderFlagBucket,
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
	driver := c.String(uploaderFlagDriver)
	endpoint := c.String(uploaderFlagEndpoint)
	accessKey := c.String(uploaderFlagAccessKey)
	accessSecret := c.String(uploaderFlagAccessSecret)
	bucketName := c.String(uploaderFlagBucket)
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
