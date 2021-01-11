package config

import (
	"google.golang.org/grpc"
	"net"
)

var (
	ConfigFile string
)

var (
	GRPCServer *grpc.Server
	GRPCListen net.Listener
)
