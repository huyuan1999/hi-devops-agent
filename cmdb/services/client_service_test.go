package services

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"testing"
)

const testGRPCServer = "127.0.0.1:8088"

func TestClientService_GetNIC(t *testing.T) {
	conn, err := grpc.Dial(testGRPCServer, grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	client := NewClientClient(conn)
	resp, err := client.GetNIC(context.Background(), &RequestClient{})
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(resp.Nic)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GRPC调用结果: ", string(data))
}
