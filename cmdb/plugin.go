package cmdb

import (
	"github.com/huyuan1999/hi-devops-agent/cmdb/services"
	"google.golang.org/grpc"
)

func Register(server *grpc.Server) {
	services.RegisterClientServer(server, new(services.ClientService))
}
