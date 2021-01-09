package config

import (
	"google.golang.org/grpc"
	"net"
)

var (
	CfgPath     string
	PluginInfo  []map[string]string
)

var (
	GRPCServer *grpc.Server
	GRPCListen net.Listener
)
