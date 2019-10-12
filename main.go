package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli"

	"uptoc/oss"
)

const (
	// uploader configs
	UPTOC_UPLOADER_OSS = "oss"

	// oss configs from os envs
	UPTOC_UPLOADER_KEYID     = "UPTOC_UPLOADER_KEYID"
	UPTOC_UPLOADER_KEYSECRET = "UPTOC_UPLOADER_KEYSECRET"

	// config from cmd flags
	UPTOC_UPLOADER  = "uploader"
	UPTOC_ENDPOINT  = "endpoint"
	UPTOC_KEYID     = "access_key"
	UPTOC_KEYSECRET = "access_secret"
	UPTOC_BUCKET    = "bucket"
)

var (
	appFlags = []cli.Flag{
		cli.StringFlag{
			Name:  UPTOC_UPLOADER,
			Usage: "specify cloud storage engine. default: oss",
			Value: "oss",
		},
		cli.StringFlag{
			Name:     UPTOC_ENDPOINT,
			Usage:    "specify endpoint of the uploader",
			Required: true,
		},
		cli.StringFlag{
			Name:     UPTOC_KEYID,
			Usage:    "specify endpoint of the uploader",
			EnvVar:   UPTOC_UPLOADER_KEYID,
			Required: true,
		},
		cli.StringFlag{
			Name:     UPTOC_KEYSECRET,
			Usage:    "specify endpoint of the uploader",
			EnvVar:   UPTOC_UPLOADER_KEYSECRET,
			Required: true,
		},
		cli.StringFlag{
			Name:     UPTOC_BUCKET,
			Usage:    "specify bucket name of the uploader",
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
	//app.Version = version.Long
	app.Flags = appFlags
	app.Action = func(c *cli.Context) {
		endpoint := c.String(UPTOC_ENDPOINT)
		keyId := c.String(UPTOC_KEYID)
		keySecret := c.String(UPTOC_KEYSECRET)
		bucketName := c.String(UPTOC_BUCKET)

		// select uploader
		var uploader Uploader
		switch c.String(UPTOC_UPLOADER) {
		case UPTOC_UPLOADER_OSS:
			ossUploader, err := oss.NewUploader(endpoint, keyId, keySecret, bucketName)
			if err != nil {
				log.Fatalln(err)
			}

			uploader = ossUploader
		}

		if err := walkAndUpload(c.Args().First(), uploader); err != nil {
			log.Fatalln(err)
		}

	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// walk and upload to the cloud storage
type Uploader interface {
	Upload(object, rawPath string) error
}

func walkAndUpload(dirPth string, uploader Uploader) error {
	return filepath.Walk(dirPth, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		fmt.Println(path)
		return uploader.Upload(strings.TrimLeft(path, dirPth+"/"), path)
	})
}
