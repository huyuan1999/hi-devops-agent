package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"github.com/huyuan1999/hi-devops-agent/apiserver/utils"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"syscall"
)

const (
	defaultNtpServer = "ntp1.aliyun.com"
	defaultZone      = "Asia/Shanghai"
)

type datetime struct {
	context *cli.Context
}

func setDatetime(context *cli.Context) {
	d := &datetime{context: context}
	d.setting()
}

func (d *datetime) setting() {
	if !d.context.Bool("set-zone") {
		return
	}

	d.setTimeZone(d.context.String("zone"))

	if d.context.Bool("set-time") {
		c := cron.New()
		_, err := c.AddFunc("* * * * *", func() {
			d.startNtpService(d.context.String("ntp-server"))
		})
		if err != nil {
			logrus.Warningf("添加 ntp 时间同步任务错误: ", err.Error())
		}
	}
}

func (d *datetime) setTimeZone(zone string) {
	if zone == "" {
		zone = defaultZone
	}
	if err := os.Setenv("TZ", zone); err != nil {
		logrus.Warningf("set TZ env error: %s", err.Error())
	}

	zoneFile := fmt.Sprintf("/usr/share/zoneinfo/%s", zone)
	if !utils.IsFile(zoneFile) {
		logrus.Warningf("zone file %s not found", zone)
		return
	}

	ln := fmt.Sprintf("ln -sf %s /etc/localtime", zone)
	cmd := exec.Command("/bin/bash", "-c", ln)
	if err := cmd.Run(); err != nil {
		logrus.Warningf("ln zone file err: %s", err.Error())
	}
}

func (d *datetime) startNtpService(server string) {
	if server == "" {
		server = defaultNtpServer
	}
	response, err := ntp.Time(server)
	if err != nil {
		logrus.Warningf("ntp client err: %s", err.Error())
		return
	}

	val := syscall.NsecToTimeval(response.UnixNano())
	if err := syscall.Settimeofday(&val); err != nil {
		logrus.Warningf("syscall set time of day err: %s", err.Error())
		return
	}
}
