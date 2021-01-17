package utils

import (
	"fmt"
	"github.com/beevik/ntp"
	"github.com/huyuan1999/hi-devops-agent/apiserver/config"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

func SetZone() {
	if config.TimeZone == "" {
		log.Warning("time zone is empty")
		return
	}
	if err := os.Setenv("TZ", config.TimeZone); err != nil {
		log.Warningf("setenv TZ err: %s", err.Error())
	}

	zone := fmt.Sprintf("/usr/share/zoneinfo/%s", config.TimeZone)

	if !IsFile(zone) {
		log.Warningf("zone file %s not found", zone)
		return
	}

	setZone := fmt.Sprintf("ln -sf %s /etc/localtime", zone)
	cmd := exec.Command("/bin/bash", "-c", setZone)
	if err := cmd.Run(); err != nil {
		log.Warningf("ln zoneinfo err: %s", err.Error())
	}
}

func SyncTime() {
	response, err := ntp.Time(config.NtpServer)
	if err != nil {
		log.Warningf("ntp err: %s", err.Error())
		return
	}
	timeval := syscall.NsecToTimeval(response.UnixNano())
	if err := syscall.Settimeofday(&timeval); err != nil {
		log.Warningf("syscall settimeofday err: %s", err.Error())
		return
	}
}
