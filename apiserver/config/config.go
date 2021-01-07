package config

import (
	"google.golang.org/grpc"
	"net"
)

var (
	CfgPath = ""
)

var (
	GRPCServer *grpc.Server
	GRPCListen net.Listener
)
