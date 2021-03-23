package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func SaveInS3(c *cli.Context, input string, logger *logrus.Entry) error {
	s3Client, err := minio.New(
		fmt.Sprintf("%s:%s", c.String("s3-server"), c.String("s3-server-port")),
		&minio.Options{Creds: credentials.NewStaticV4(
			c.String("access-key"),
			c.String("secret-key"),
			"",
		)},
	)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in getting minio client %s", err),
			2,
		)
	}
	info, err := s3Client.FPutObject(
		context.Background(),
		c.String("s3-bucket"),
		c.String("upload-path"),
		input,
		minio.PutObjectOptions{ContentType: "application/text"},
	)
	if err != nil {
		fmt.Errorf("unable to upload file %s", err)
	}
	logger.Infof(
		"save file %s to s3 storage with etag %s and version %s",
		input, info.ETag, info.VersionID,
	)
	return nil
}
