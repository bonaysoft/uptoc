package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"

	"uptoc/core"
	"uptoc/uploader"
	"uptoc/version"
)

const (
	// core configs
	UPTOC_UPLOADER_OSS = "oss"

	// uploader configs from os envs
	UPTOC_UPLOADER_KEYID     = "UPTOC_UPLOADER_KEYID"
	UPTOC_UPLOADER_KEYSECRET = "UPTOC_UPLOADER_KEYSECRET"

	// config from cmd flags
	UPTOC_UPLOADER  = "core"
	UPTOC_ENDPOINT  = "endpoint"
	UPTOC_KEYID     = "access_key"
	UPTOC_KEYSECRET = "access_secret"
	UPTOC_BUCKET    = "bucket"
)

var (
	appFlags = []cli.Flag{
		cli.StringFlag{
			Name:  UPTOC_UPLOADER,
			Usage: "specify cloud storage engine. default: uploader",
			Value: UPTOC_UPLOADER_OSS,
		},
		cli.StringFlag{
			Name:     UPTOC_ENDPOINT,
			Usage:    "specify endpoint of the core",
			Required: true,
		},
		cli.StringFlag{
			Name:     UPTOC_KEYID,
			Usage:    "specify endpoint of the core",
			EnvVar:   UPTOC_UPLOADER_KEYID,
			Required: true,
		},
		cli.StringFlag{
			Name:     UPTOC_KEYSECRET,
			Usage:    "specify endpoint of the core",
			EnvVar:   UPTOC_UPLOADER_KEYSECRET,
			Required: true,
		},
		cli.StringFlag{
			Name:     UPTOC_BUCKET,
			Usage:    "specify bucket name of the core",
			Required: true,
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "uptoc"
	app.Usage = "A cli tool to upload the dist file for the cloud engine."
	app.Copyright = "(c) 2019 uptoc.saltbo.cn"
	app.Compiled = time.Now()
	app.Version = version.Long
	app.Flags = appFlags
	app.Action = appAction
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func appAction(c *cli.Context) {
	driver := c.String(UPTOC_UPLOADER)
	endpoint := c.String(UPTOC_ENDPOINT)
	keyId := c.String(UPTOC_KEYID)
	keySecret := c.String(UPTOC_KEYSECRET)
	bucketName := c.String(UPTOC_BUCKET)
	u, err := uploader.New(driver, endpoint, keyId, keySecret, bucketName)
	if err != nil {
		log.Fatalln(err)
	}

	dirPath := c.Args().First()
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}

	e := core.NewEngine(u)
	if err := e.LoadAndCompareObjects(dirPath); err != nil {
		log.Fatalln(err)
	}

	if err := e.Sync(); err != nil {
		log.Fatalln(err)
	}
}
