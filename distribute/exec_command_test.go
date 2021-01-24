package distribute

import (
	"context"
	"github.com/huyuan1999/hi-devops-agent/distribute/services"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	options := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(115343360),
		grpc.MaxSendMsgSize(115343360),
	}
	server := grpc.NewServer(options...)
	Register(server)
	listen, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalln(err)
	}

	if err := server.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}

func TestClient(t *testing.T) {
	// 设置客户端最大接收大小为 100M
	opts := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024 * 1024 * 100))

	conn, err := grpc.Dial("192.168.3.10:8089", grpc.WithInsecure(), opts)
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = conn.Close() }()

	client := services.NewExecClient(conn)
	resp, err := client.Output(context.Background(), &services.Command{Cmd: "ls", Args: []string{"-h", "-l", "-a", "/root/"}})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("服务器返回的值: \n", resp.Out)
}
