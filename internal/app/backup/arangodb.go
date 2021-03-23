package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/dictybase-docker/k8s-custodian/internal/logger"
	"github.com/dictybase-docker/k8s-custodian/internal/storage"
	"github.com/mholt/archiver/v3"
	"github.com/urfave/cli"
)

func RunArangoBackup(c *cli.Context) error {
	dumpCmd, err := exec.LookPath("arangodump")
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("error in finding arangodump executables %s", err),
			2,
		)
	}
	// dump the database
	dumpDir, err := outDir(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	defer os.RemoveAll(dumpDir)
	args := []string{
		"--server.endpoint",
		fmt.Sprintf("ssl://%s:%s", c.String("arangodb-host"), c.String("arangodb-port")),
		"--server.username", c.String("arangodb-user"),
		"--server.password", c.String("arangodb-pass"),
		"--server.database", c.String("arangodb-database"),
		"--compress-output", "--dump-dependencies",
		"--include-system-collection",
		"--threads", "4",
		"--output-directory", dumpDir,
	}
	cmd := exec.Command(dumpCmd, args...)
	stdStderr, err := cmd.CombinedOutput()
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf(
				"error in running command dump command with args %s %s %s",
				args, string(stdStderr), err,
			),
			2,
		)
	}
	aFile, err := archiveDir(dumpDir)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	err = storage.SaveInS3(c, aFile, logger.GetLogger(c))
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	return nil
}

func archiveDir(dir string) (string, error) {
	aDir, err := ioutil.TempDir(os.TempDir(), "archive-*")
	if err != nil {
		return aDir,
			fmt.Errorf("error in creating a temp dir for archive %s", err)
	}
	aFile := filepath.Join(
		aDir,
		fmt.Sprintf("arangobackup-%s.tar", time.Now().Format("01-02-2006")),
	)
	return aFile, archiver.Archive([]string{dir}, aFile)
}

func outDir(c *cli.Context) (string, error) {
	dirPrefix := fmt.Sprintf("%s-%s", c.String("arangodb-database"), time.Now().Format("01-02-2006"))
	parentDir := os.TempDir()
	dumpDir, err := ioutil.TempDir(parentDir, fmt.Sprintf("%s-*", dirPrefix))
	if err != nil {
		return dumpDir,
			fmt.Errorf("error in creating a temp dir with prefix %s %s", dirPrefix, err)
	}
	return dumpDir, nil
}
