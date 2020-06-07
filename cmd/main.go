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
			Name:     uploaderFlagRegion,
			Usage:    "specify region of the cloud platform",
			Required: true,
		},
		cli.StringFlag{
			Name:     uploaderFlagAccessKey,
			Usage:    "specify key id of the cloud platform",
			EnvVar:   uploaderEnvAccessKey,
			Required: true,
		},
		cli.StringFlag{
			Name:     uploaderFlagSecretKey,
			Usage:    "specify key secret of the cloud platform",
			EnvVar:   uploaderEnvSecretKey,
			Required: true,
		},
		cli.StringFlag{
			Name:     uploaderFlagBucket,
			Usage:    "specify bucket name of the cloud platform",
			Required: true,
		},
		cli.StringFlag{
			Name:     uploaderFlagExclude,
			Usage:    "specify exclude the given comma separated directories (example: --exclude=.cache,test)",
			Required: false,
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
	var excludePaths []string
	ak := c.String(uploaderFlagAccessKey)
	sk := c.String(uploaderFlagSecretKey)
	driver := c.String(uploaderFlagDriver)
	region := c.String(uploaderFlagRegion)
	bucket := c.String(uploaderFlagBucket)
	exclude := c.String(uploaderFlagExclude)
	if exclude != "" {
		excludePaths = strings.Split(exclude, ",")
	}

	dirPath := c.Args().First()
	log.Printf("driver: %s\n", driver)
	log.Printf("region: %s\n", region)
	log.Printf("bucket: %s\n", bucket)
	log.Printf("exclude: %s\n", excludePaths)
	log.Printf("dirPath: %s\n", dirPath)
	uploadDriver, err := uploader.New(driver, region, ak, sk, bucket)
	if err != nil {
		log.Fatalln(err)
	}

	e := core.NewEngine(uploadDriver)
	if err := e.LoadAndCompareObjects(dirPath, excludePaths...); err != nil {
		log.Fatalln(err)
	}

	if err := e.Sync(); err != nil {
		log.Fatalln(err)
	}
}
