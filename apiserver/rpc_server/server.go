package rpc_server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/huyuan1999/hi-devops-agent/apiserver/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
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

// grpc 双向证书验证
func (rpc *RPCServer) authentication() (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.CAPem)
	if err != nil {
		log.Fatalln(err)
	}
	certPool.AppendCertsFromPEM(ca)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	return creds, err
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

	creds, err := rpc.authentication()
	if err != nil {
		return nil, nil, err
	}

	server := grpc.NewServer(grpc.Creds(creds))
	return server, listen, nil
}
