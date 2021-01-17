package main

import (
	"encoding/json"
	"flag"
	"github.com/huyuan1999/hi-devops-agent/apiserver/config"
	"github.com/huyuan1999/hi-devops-agent/apiserver/rpc_server"
	"github.com/huyuan1999/hi-devops-agent/apiserver/utils"
	"github.com/huyuan1999/hi-devops-agent/cmdb"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

type initData struct {
	OnlyId    string `json:"olny_id"`
	NtpServer string `json:"ntp_server"`
	TimeZone  string `json:"time_zone"`
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func mkdir(name string) {
	if !utils.IsDir(name) {
		if err := os.Mkdir(name, 0755); err != nil {
			log.Fatalln("makedir error: ", err.Error())
		}
	}
}

func init() {
	log.SetReportCaller(true)

	var cfgPath string
	var cfgPrint bool

	flag.StringVar(&cfgPath, "config", "/etc/devops/hi-devops-agent.ini", "指定配置文件")
	flag.BoolVar(&cfgPrint, "print", false, "打印配置文件示例")
	flag.Parse()

	if cfgPrint {
		cfgTemplate()
		os.Exit(0)
	}

	analysisConfig(cfgPath)

	if !utils.Exists(config.CfgWorkdir) {
		mkdir(config.CfgWorkdir)
	}

	// 切换到工作目录
	err := os.Chdir(config.CfgWorkdir)
	Panic(err)

	mkdir("tls")
	config.ConfigFile = cfgPath
}

func setTime() {
	if config.CfgSetZone {
		utils.SetZone()
	}
	if config.CfgSetTime {
		utils.SyncTime()
		go func() {
			for {
				utils.SyncTime()
				time.Sleep(time.Second * 60)
			}
		}()
	}
}

func start() {
	if utils.IsFile(config.InitializationFile) {
		// 加载 .init.json
		d := &initData{}
		file, err := ioutil.ReadFile(config.InitializationFile)
		Panic(err)
		err = json.Unmarshal(file, d)
		Panic(err)
		config.OnlyId = d.OnlyId
		config.TimeZone = d.TimeZone
		config.NtpServer = d.NtpServer
	} else {
		// 检查是否第一次启动 agent, 是否需要初始化
		// 注意: 同一个 agent 节点不能重复初始化(服务器端将会比较 HostID, 如果发现 HostID 已存在则会拒绝相关操作)
		initialization := NewInitialization()
		err := initialization.Register()
		Panic(err)
	}

	setTime()

	// 初始化 grpc 服务
	rpc := rpc_server.NewRPC("tcp", config.CfgAddress)
	server, listen, err := rpc.Listen()
	Panic(err)

	// 加载插件
	cmdb.Register(server)

	err = server.Serve(listen)
	Panic(err)
}

func main() {
	if config.CfgDaemon {
		Daemon(start, config.CfgPidFile, config.CfgLogFile, config.CfgWorkdir)
	} else {
		start()
	}
}
