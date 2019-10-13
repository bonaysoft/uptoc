package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"

	"uptoc/oss"
	"uptoc/utils"
	"uptoc/version"
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

// walk and upload to the cloud storage
type Uploader interface {
	ListObjects() ([]string, error)
	Upload(object, rawPath string) error
	Delete(object string) error
}

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

	if err := sync2Cloud(uploader, c.Args().First()); err != nil {
		log.Fatalln(err)
	}
}

func sync2Cloud(uploader Uploader, dirPath string) error {
	if !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}

	log.Printf("find the local objects...")
	files, err := utils.ListLocalFiles(dirPath)
	if err != nil {
		return err
	}

	// transfer the localFiles to localObjects
	localObjects := make([]string, 0, len(files))
	for _, filePath := range files {
		object := strings.TrimPrefix(filePath, dirPath)
		localObjects = append(localObjects, object)
	}

	log.Printf("find the remote objects...")
	remoteObjects, err := uploader.ListObjects()
	if err != nil {
		return err
	}

	log.Printf("compare and find the deleted files...")
	deletedObjects := make([]string, 0)
	for _, obj := range remoteObjects {
		if !utils.StrInSlice(obj, localObjects) {
			deletedObjects = append(deletedObjects, obj)
		}
	}

	if err := uploadObjects(uploader, dirPath, localObjects); err != nil {
		return err
	}

	cleanDeletedObjects(uploader, deletedObjects)
	return nil
}

func uploadObjects(uploader Uploader, dirPath string, objects []string) error {
	log.Printf("found %d local files, uploading...", len(objects))
	for _, object := range objects {
		if err := uploader.Upload(object, dirPath+object); err != nil {
			return err
		}
	}

	log.Printf("upload files done.")
	return nil
}

func cleanDeletedObjects(uploader Uploader, objects []string) {
	log.Printf("found %d deleted files, cleaning...", len(objects))
	for _, obj := range objects {
		err := uploader.Delete(obj)
		if err != nil {
			log.Printf("remove the file %s failed: %s", err)
		}
	}

	log.Printf("clean deleted files done.")
}
