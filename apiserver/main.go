package main

import (
	"github.com/huyuan1999/hi-devops-agent/apiserver/config"
	"github.com/huyuan1999/hi-devops-agent/apiserver/rpc_server"
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
	// cmdb.Register(config.GRPCServer)
	// 编译插件 SVCNAME=$1 SVCVER=$2 TIMESTAMP=`date '+%Y%m%d_%H%M%S'` go build -v -buildmode=plugin --ldflags="-pluginpath=${SVCNAME}_${TIMESTAMP}" -o ${SVCNAME}_${SVCVER}.so ${SVCNAME}
	// 使用配置文件加载插件
	/*
		[plugin:cmdb]
		name = "cmdb"
		path = "/usr/local/hi-devops-agent/plugins/cmdb_1.0.0.so"
	*/

	// 注册插件
	NewRegister().Load()

	if err := config.GRPCServer.Serve(config.GRPCListen); err != nil {
		panic(err)
	}
}
