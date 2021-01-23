package main

import (
	"fmt"
	"github.com/sevlyar/go-daemon"
	"github.com/urfave/cli/v2"
	"syscall"
)

var stopCommand = &cli.Command{
	Name:  "stop",
	Usage: "停止以 daemon 方式启动的 agent",
	Action: func(context *cli.Context) error {
		pid, err := daemon.ReadPidFile("/run/devops.pid")
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			fmt.Println(err.Error())
		}
		return nil
	},
}
