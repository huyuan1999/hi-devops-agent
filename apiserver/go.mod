module github.com/huyuan1999/hi-devops-agent/apiserver

go 1.15

require (
	github.com/beevik/ntp v0.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/huyuan1999/hi-devops-agent/cmdb v0.0.0-20210111080328-e9247570bdfe
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.0
	github.com/sevlyar/go-daemon v0.1.5
	github.com/shirou/gopsutil/v3 v3.20.12
	github.com/sirupsen/logrus v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/sys v0.0.0-20210113000019-eaf3bda374d2 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210111234610-22ae2b108f89 // indirect
	google.golang.org/grpc v1.34.1
	google.golang.org/protobuf v1.25.0
)

replace github.com/huyuan1999/hi-devops-agent/cmdb v0.0.0-20210111080328-e9247570bdfe => ../cmdb
