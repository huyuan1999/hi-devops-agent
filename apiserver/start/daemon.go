package start

import (
	"fmt"
	"github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
	"os"
)

func Daemon(process func(), pidFileName, logFileName, workDir string) {
	ctx := &daemon.Context{
		PidFileName: pidFileName,
		PidFilePerm: 0644,
		LogFileName: logFileName,
		LogFilePerm: 0640,
		WorkDir:     workDir,
		Umask:       027,
		Args:        os.Args,
	}

	child, err := ctx.Reborn()
	if err != nil {
		logrus.Error(fmt.Sprintf("fork daemon error: %v", err))
		return
	}

	if child != nil {
		return
	}

	defer func() { _ = ctx.Release() }()
	logrus.Info("daemon started [ OK ]")

	process()

	logrus.Info("daemon terminated")
}
