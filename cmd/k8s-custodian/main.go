package main

import (
	"github.com/dictybase-docker/k8s-custodian/internal/app/backup"
	"github.com/dictybase-docker/k8s-custodian/internal/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "custodian"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-format",
			Usage: "format of the logging out, either of json or text.",
			Value: "json",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "log level for the application",
			Value: "error",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "arangoback-to-minioS3",
			Usage:  "backup arangodb database to minio s3 storage",
			Flags:  getArangoBackupFlags(),
			Action: backup.RunArangoBackup,
		},
	}
}

func getArangoBackupFlags() []cli.Flag {
	var f []cli.Flag
	f = append(f, cmd.ArangodbBackupCmd()...)
	return append(f, cmd.MinioS3StorageCmd()...)
}
