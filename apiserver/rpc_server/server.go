package rpc_server

import (
	"errors"
	"google.golang.org/grpc"
	"net"
	"strings"
)

func NewRPC(protocol, address string) *RPCServer {
	return &RPCServer{
		Protocol: protocol,
		Address:  address,
	}
}

type RPCServer struct {
	Protocol string
	Address  string
}

func (rpc *RPCServer) Listen() (*grpc.Server, net.Listener, error) {
	switch strings.ToUpper(rpc.Protocol) {
	case "TCP":
		return rpc.tcp()
	default:
		return nil, nil, errors.New("unexpected address type")
	}
}

func (rpc *RPCServer) tcp() (*grpc.Server, net.Listener, error) {
	listen, err := net.Listen("tcp", rpc.Address)
	if err != nil {
		return nil, nil, err
	}

	server := grpc.NewServer()
	return server, listen, nil
}
