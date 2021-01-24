package distribute

import (
	"github.com/huyuan1999/hi-devops-agent/distribute/services"
	"google.golang.org/grpc"
)


func Register(server *grpc.Server) {
	services.RegisterExecServer(server, new(services.ExecCommandService))
	services.RegisterDistributeServer(server, new(services.FileService))
}
