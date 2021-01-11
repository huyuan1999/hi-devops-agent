package command

import (
	"github.com/huyuan1999/hi-devops-agent/cmdb/utils"
	"regexp"
)

type Smart struct {
}

func NewSmart() (*Smart, error) {
	if err := utils.Which("smartctl"); err != nil {
		return nil, err
	} else {
		return &Smart{}, nil
	}
}

func (s *Smart) re(device, compile string) string {
	result := utils.Shell("smartctl", "-i", device)
	reg := regexp.MustCompile(compile)
	matched := reg.FindStringSubmatch(result.Stdout)
	if len(matched) == 2 {
		return utils.Trim(matched[1])
	}
	return ""
}

func (s *Smart) Info(device string) string {
	result := utils.Shell("smartctl", "-i", device)
	return result.Stdout
}

func (s *Smart) ModelFamily(device string) string {
	family := s.re(device, "Model\\s+Family:\\s+(.*)")
	if family != "" {
		return family
	}
	return s.re(device, "Vendor:\\s+(.*)")
}

func (s *Smart) DeviceModel(device string) string {
	model := s.re(device, "Device\\s+Model:\\s+(.*)")
	if model != "" {
		return model
	}
	return s.re(device, "Product:\\s+(.*)")
}

func (s *Smart) SerialNumber(device string) string {
	sn := s.re(device, "Serial\\s+Number:\\s+(.*)")
	if sn != "" {
		return sn
	}
	return s.re(device, "Serial\\s+number:\\s+(.*)")
}

func (s *Smart) FormFactor(device string) string {
	return s.re(device, "Form\\s+Factor:\\s+(.*)")
}

func (s *Smart) RotationRate(device string) string {
	return s.re(device, "Rotation\\s+Rate:\\s+(.*)")
}

func (s *Smart) UserCapacity(device string) string {
	return s.re(device, "User\\s+Capacity:\\s+(.*)")
}
