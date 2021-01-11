package linux

import (
	"fmt"
	"github.com/huyuan1999/hi-devops-agent/cmdb/models"
	"github.com/huyuan1999/hi-devops-agent/cmdb/utils"
	"github.com/shirou/gopsutil/v3/host"
)

type System struct {
	models.System
	info *host.InfoStat
}

func NewSystem() (*System, error) {
	system := &System{}
	info, err := host.Info()
	if err != nil {
		return nil, err
	}
	system.info = info
	utils.Call(system)
	return system, nil
}

func (s *System) GetHostName() {
	s.HostName = s.info.Hostname
}

func (s *System) GetKernel() {
	s.Kernel = s.info.KernelVersion
}

func (s *System) GetRelease() {
	s.Release = fmt.Sprintf("%s %s", s.info.Platform, s.info.PlatformVersion)
}
