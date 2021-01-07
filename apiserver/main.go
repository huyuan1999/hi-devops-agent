package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/huyuan1999/hi-devops-agent/apiserver/config"
	"github.com/huyuan1999/hi-devops-agent/apiserver/plugins"
	"github.com/huyuan1999/hi-devops-agent/apiserver/rpc_server"
	"log"
	"os"
)

func init() {
	rpc := rpc_server.NewRPC("tcp", ":8088")
	server, listen, err := rpc.Listen()
	if err != nil {
		panic(err)
	}
	config.GRPCListen = listen
	config.GRPCServer = server
}

func main() {
	// 加载组件
	// 包括 cmdb, 批量命令, webssh, 配置管理, 发布平台 等
	// cmdb.Init(config.GRPCServer)
	// 使用配置文件加载插件
	/*
	[plugin:cmdb]
	name = "cmdb"
	path = "/usr/local/hi-devops-agent/plugins/cmdb"
	*/
	//if err := config.GRPCServer.Serve(config.GRPCListen); err != nil {
	//	panic(err)
	//}

	pl := plugins.Plugin{
		LoggerOutput: os.Stdout,
		LoggerLevel:  hclog.Info,
	}

	log.Println("调用插件: ", pl.Load("cmdb"))
}
