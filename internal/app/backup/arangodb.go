package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/dictybase-docker/k8s-custodian/internal/logger"
	"github.com/dictybase-docker/k8s-custodian/internal/storage"
	"github.com/mholt/archiver/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func ArangoBackupToMinioS3(c *cli.Context) error {
	logger := log.GetLogger(c)
	// dump the database
	dumpDir, err := arangoDump(c, logger)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	defer os.RemoveAll(dumpDir)
	// create tar archive
	aDir, aFile, err := archiveDir(dumpDir, logger)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	defer os.RemoveAll(aDir)
	// upload to minio s3 storage
	err = storage.SaveInS3(c, aFile, logger)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	return nil
}

func arangoDump(c *cli.Context, logger *logrus.Entry) (string, error) {
	dumpCmd, err := exec.LookPath("arangodump")
	if err != nil {
		return "", fmt.Errorf("error in finding arangodump executables %s", err)
	}
	dumpDir, err := outDir(c)
	if err != nil {
		return dumpDir, err
	}
	logger.Debugf("going to dump database in %s", dumpDir)
	args := []string{
		"--server.endpoint",
		fmt.Sprintf("ssl://%s:%s", c.String("arangodb-host"), c.String("arangodb-port")),
		"--server.username", c.String("arangodb-user"),
		"--server.password", c.String("arangodb-pass"),
		"--server.database", c.String("arangodb-database"),
		"--compress-output", "--dump-dependencies",
		"--include-system-collections",
		"--threads", "4",
		"--output-directory", dumpDir,
	}
	cmd := exec.Command(dumpCmd, args...)
	logger.Debugf("going to run dump command %s", args)
	stdStderr, err := cmd.CombinedOutput()
	if err != nil {
		return dumpDir, fmt.Errorf(
			"error in running command dump command with args %s %s %s",
			args, string(stdStderr), err,
		)
	}
	logger.Debugf("dump output %s", string(stdStderr))
	logger.Infof("dumped the database %s", c.String("arangodb-database"))
	return dumpDir, nil
}

func archiveDir(dir string, logger *logrus.Entry) (string, string, error) {
	aDir, err := os.MkdirTemp(os.TempDir(), "*-archive")
	if err != nil {
		return aDir, "",
			fmt.Errorf("error in creating a temp dir for archive %s", err)
	}
	logger.Debugf("going to create archive in %s folder", aDir)
	aFile := filepath.Join(
		aDir,
		fmt.Sprintf("arangobackup-%s.tar", time.Now().Format("01-02-2006")),
	)
	logger.Infof("created archive %s for dir %s", aFile, dir)
	return aDir, aFile, archiver.Archive([]string{dir}, aFile)
}

func outDir(c *cli.Context) (string, error) {
	dirSuffix := fmt.Sprintf("%s-%s", c.String("arangodb-database"), time.Now().Format("01-02-2006"))
	dumpDir, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("*-%s", dirSuffix))
	if err != nil {
		return dumpDir,
			fmt.Errorf("error in creating a temp dir with prefix %s %s", dirSuffix, err)
	}
	return dumpDir, nil
}
