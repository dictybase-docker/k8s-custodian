package cmd

import (
	arangoflag "github.com/dictyBase/arangomanager/command/flag"
	"github.com/urfave/cli"
)

func ArangodbBackupCmd() []cli.Flag {
	return arangoflag.ArangodbFlags()
}
