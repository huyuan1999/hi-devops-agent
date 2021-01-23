package start

import (
	"github.com/huyuan1999/hi-devops-agent/apiserver/utils"
	"github.com/huyuan1999/hi-devops-agent/cmdb"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var StartCommand = &cli.Command{
	Name:  "start",
	Usage: "启动 agent",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "daemon",
			Aliases: []string{"d"},
			Value:   false,
			Usage:   "以守护进程方式启动",
		},
		&cli.StringFlag{
			Name:   "pid",
			Value:  "/run/devops.pid",
			Usage:  "进程 pid 存放位置, 此选项只有在 daemon=true 时生效",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:  "listen",
			Value: "0.0.0.0:8088",
			Usage: "agent 监听地址",
		},
	},
	Action: func(context *cli.Context) error {
		if !utils.IsFile(".init.json") {
			logrus.Fatalln("请先执行 init 操作")
		}
		if context.Bool("daemon") {
			logFile := context.String("log")
			pidFile := context.String("pid")
			workDir := context.String("work-dir")
			Daemon(func() { start(context) }, pidFile, logFile, workDir)
		} else {
			start(context)
		}
		return nil
	},
}

func start(context *cli.Context) {
	listen := context.String("listen")
	rpc := &RPCServer{
		Protocol: "tcp",
		Address:  listen,
	}

	cert := context.String("cert-pem")
	key := context.String("cert-key")
	ca := context.String("ca")

	server, listener, err := rpc.Listen(key, cert, ca)
	if err != nil {
		utils.FatalError(err)
	}

	// 加载插件
	cmdb.Register(server)

	if err = server.Serve(listener); err != nil {
		utils.FatalError(err)
	}
}
