package main

import (
	"flag"
	"github.com/huyuan1999/hi-devops-agent/apiserver/config"
	"github.com/huyuan1999/hi-devops-agent/apiserver/rpc_server"
)

func init() {
	rpc := rpc_server.NewRPC("tcp", ":8088")
	server, listen, err := rpc.Listen()
	if err != nil {
		panic(err)
	}

	var configFile string
	flag.StringVar(&configFile, "config", "./hi-devops-agent.ini", "指定配置文件")
	defer flag.Parse()

	config.ConfigFile = configFile
	config.GRPCListen = listen
	config.GRPCServer = server
}

func main() {
	// 加载组件
	// 包括 cmdb, 批量命令, webssh, 配置管理, 发布平台 等
	// cmdb.Register(config.GRPCServer)
	//cmdb.Register(config.GRPCServer)

	if err := config.GRPCServer.Serve(config.GRPCListen); err != nil {
		panic(err)
	}
}
