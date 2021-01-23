package initialization

import (
	"fmt"
	"github.com/huyuan1999/hi-devops-agent/apiserver/utils"
	"github.com/urfave/cli/v2"
)

var InitializationCommand = &cli.Command{
	Name:  "init",
	Usage: "初始化节点",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "public-ip",
			Required: true,
			Usage:    "当前节点在集群中的标识 IP, 通常应该是唯一的",
		},
		&cli.StringFlag{
			Name:     "server",
			Required: true,
			Usage:    "服务器地址, [http|https]://address:port/",
		},
	},
	Action: func(context *cli.Context) error {
		err := initialization(context)
		if err != nil {
			utils.FatalError(err)
		} else {
			fmt.Println("初始化完成, 请运行 start 启动 agent")
		}
		return nil
	},
}
