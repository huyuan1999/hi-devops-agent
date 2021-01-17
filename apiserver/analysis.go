package main

import (
	"fmt"
	"github.com/huyuan1999/hi-devops-agent/apiserver/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func cfgTemplate() {
	fmt.Printf(`[global]
# 服务器地址
server = devops.huyuan.io:80

# agent 工作目录
workdir = "/opt/hi-devops-agent/"

# agent grpc 监听地址
address = "0.0.0.0:8088"

# 是否以守护进程方式运行, 如果为 true 则需要指定 pidfile 和 logfile
daemon = false
pidfile = ""
logfile = ""

# 是否自动校准系统时间(以 server 系统时间为标准)
settime = true

# 是否自动校准系统时区(以 server 系统时区为标准)
setzone = true

# 可选: 标识服务器的 ip, 如果指定了 server webssh 将使用这个地址连接此服务器
public_ip = "127.0.0.1"
`)
}

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func analysisConfig(path string) {
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatalln(err)
	}

	config.CfgServer = cfg.Section("global").Key("server").String()
	config.CfgWorkdir = cfg.Section("global").Key("workdir").String()
	config.CfgAddress = cfg.Section("global").Key("address").String()
	config.CfgPidFile = cfg.Section("global").Key("pidfile").String()
	config.CfgLogFile = cfg.Section("global").Key("logfile").String()
	config.CfgPublicIP = cfg.Section("global").Key("public_ip").String()

	daemon, err := cfg.Section("global").Key("daemon").Bool()
	Fatal(err)
	settime, err := cfg.Section("global").Key("settime").Bool()
	Fatal(err)
	setzone, err := cfg.Section("global").Key("setzone").Bool()
	Fatal(err)

	config.CfgDaemon = daemon
	config.CfgSetTime = settime
	config.CfgSetZone = setzone
}
