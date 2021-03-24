package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
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
	if err := bucketConfiguration(c, s3Client); err != nil {
		return cli.NewExitError(err.Error(), 2)
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

func bucketConfiguration(c *cli.Context, client *minio.Client) error {
	bucket := c.String("s3-bucket")
	if err := findOrCreateBucket(client, bucket); err != nil {
		return err
	}
	return manageBucketLifecycle(client, bucket, c.Int("expiration"))
}

func manageBucketLifecycle(client *minio.Client, bucket string, expiration int) error {
	exConfig, err := client.GetBucketLifecycle(context.Background(), bucket)
	if err != nil {
		return err
	}
	if exConfig != nil {
		if int(exConfig.Rules[0].Expiration.Days) == expiration {
			return nil
		}
	}
	return configureBucketLifecycle(client, bucket, expiration)
}

func configureBucketLifecycle(client *minio.Client, bucket string, expiration int) error {
	config := lifecycle.NewConfiguration()
	config.Rules = []lifecycle.Rule{
		{
			ID:     "expire-backup-bucket",
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: lifecycle.ExpirationDays(expiration),
			},
		},
	}
	err := client.SetBucketLifecycle(context.Background(), bucket, config)
	if err != nil {
		return fmt.Errorf("error in setting bucket lifecycle %s", err)
	}
	return nil
}

func findOrCreateBucket(client *minio.Client, bucket string) error {
	ok, err := client.BucketExists(context.Background(), bucket)
	if err != nil {
		return fmt.Errorf("error in finding bucket %s", err)
	}
	if ok {
		return nil
	}
	err = client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("error in creating bucket %s", err)
	}
	return client.EnableVersioning(context.Background(), bucket)
}
