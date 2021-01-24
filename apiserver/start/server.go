package start

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/huyuan1999/hi-devops-agent/apiserver/error_type"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
	"strings"
)

type RPCServer struct {
	Protocol string
	Address  string
}

// grpc 双向证书验证
func (rpc *RPCServer) authentication(key, cert, ca string) (credentials.TransportCredentials, error) {
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, errors.Wrap(err, error_type.UnknownError)
	}

	certPool := x509.NewCertPool()
	CA, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, errors.Wrap(err, error_type.IOError)
	}
	certPool.AppendCertsFromPEM(CA)
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	return creds, nil
}

func (rpc *RPCServer) Listen(key, cert, ca string) (*grpc.Server, net.Listener, error) {
	creds, err := rpc.authentication(key, cert, ca)
	if err != nil {
		return nil, nil, err
	}
	switch strings.ToUpper(rpc.Protocol) {
	case "TCP":
		return rpc.tcp(creds)
	default:
		return nil, nil, errors.Wrap(errors.New("unexpected address type"), error_type.NotImplementedError)
	}
}

func (rpc *RPCServer) tcp(creds credentials.TransportCredentials) (*grpc.Server, net.Listener, error) {
	listen, err := net.Listen("tcp", rpc.Address)
	if err != nil {
		return nil, nil, errors.Wrap(err, error_type.ListenerError)
	}

	// 设置服务器端启用证书双向验证
	// 设置最大接收和发送大小为 115343360(110M)
	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(115343360),
		grpc.MaxSendMsgSize(115343360),
		grpc.Creds(creds),
	}

	server := grpc.NewServer(opts...)
	return server, listen, nil
}
