package config

import (
	"google.golang.org/grpc"
	"net"
)

var (
	GRPCServer *grpc.Server
	GRPCListen net.Listener
)
