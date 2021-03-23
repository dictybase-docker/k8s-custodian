package cmd

import (
	"github.com/urfave/cli"
)

func MinioS3StorageCmd() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "s3-server",
			Usage:  "S3 server endpoint",
			Value:  "minio",
			EnvVar: "MINIO_SERVICE_HOST",
		},
		cli.StringFlag{
			Name:   "s3-server-port",
			Usage:  "S3 server port",
			EnvVar: "MINIO_SERVICE_PORT",
		},
		cli.StringFlag{
			Name:     "s3-bucket",
			Usage:    "S3 bucket where the data will be uploaded",
			Required: true,
		},
		cli.StringFlag{
			Name:     "access-key, akey",
			Usage:    "access key for S3 server, required based on command run",
			Required: true,
		},
		cli.StringFlag{
			Name:     "secret-key, skey",
			Usage:    "secret key for S3 server, required based on command run",
			Required: true,
		},
		cli.StringFlag{
			Name:     "upload-path,p",
			Usage:    "full upload path inside the bucket",
			Required: true,
		},
	}
}
