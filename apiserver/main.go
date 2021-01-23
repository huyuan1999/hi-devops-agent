package main

import (
	"encoding/json"
	"github.com/huyuan1999/hi-devops-agent/apiserver/initialization"
	"github.com/huyuan1999/hi-devops-agent/apiserver/start"
	"github.com/huyuan1999/hi-devops-agent/apiserver/utils"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
)

func mkdir(name string) {
	if !utils.IsDir(name) {
		if err := os.Mkdir(name, 0755); err != nil {
			log.Fatalln("make dir error: ", err.Error())
		}
	}
}

func init() {
	logrus.SetReportCaller(true)
}

func main() {
	app := cli.NewApp()
	app.Name = "hi devops agent"
	app.Usage = ""
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "work-dir",
			Value: "/opt/hi-devops-agent/",
			Usage: "程序工作目录",
		},
		&cli.StringFlag{
			Name:  "cert-pem",
			Value: "tls/server.pem",
			Usage: "服务器证书公钥存放路径",
		},
		&cli.StringFlag{
			Name:  "cert-key",
			Value: "tls/server.key",
			Usage: "服务器证书私钥存放路径",
		},
		&cli.StringFlag{
			Name:  "ca",
			Value: "tls/ca.pem",
			Usage: "CA 证书公钥存放路径",
		},
		&cli.BoolFlag{
			Name:  "set-time",
			Value: true,
			Usage: "是否允许程序自动同步时间",
		},
		&cli.BoolFlag{
			Name:  "set-zone",
			Value: true,
			Usage: "是否允许程序设置时区",
		},
		&cli.StringFlag{
			Name:  "log",
			Value: "devops.log",
			Usage: "日志文件存放位置",
		},
		&cli.StringFlag{
			Name:  "log-format",
			Value: "text",
			Usage: "日志文件格式([ text | json ])",
		},
		&cli.StringFlag{
			Name:  "log-level",
			Value: "info",
			Usage: "日志级别([ panic | fatal | error | warn | info | debug | trace ])",
		},
		&cli.StringFlag{
			Name:   "cluster-server",
			Usage:  "服务器地址",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "only-id",
			Usage:  "唯一ID",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "zone",
			Usage:  "时区",
			Hidden: true,
		},
		&cli.StringFlag{
			Name:   "ntp-server",
			Usage:  "ntp server",
			Hidden: true,
		},
	}
	app.Commands = []*cli.Command{
		initialization.InitializationCommand,
		start.StartCommand,
		stopCommand,
	}

	app.Before = func(context *cli.Context) error {
		workDir := context.String("work-dir")
		chdir(workDir)
		mkdir("tls")

		SetLog(context)

		if utils.IsFile(".init.json") {
			loadInitFile(context)
			setDatetime(context)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln(err)
	}
}

func SetLog(context *cli.Context) {
	if context.String("log-format") == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	switch context.String("log-level") {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func chdir(workDir string) {
	if !utils.IsDir(workDir) {
		if err := os.MkdirAll(workDir, 0755); err != nil {
			logrus.Fatalln(err)
		}
	}
	if err := os.Chdir(workDir); err != nil {
		logrus.Fatalln(err)
	}
}

func loadInitFile(context *cli.Context) {
	load := make(map[string]string)
	data, err := ioutil.ReadFile(".init.json")
	if err != nil {
		logrus.Fatalln(err)
	}

	if err := json.Unmarshal(data, &load); err != nil {
		logrus.Fatalln(err)
	}

	set(context, "cluster-server", load["server"])
	set(context, "only-id", load["only_id"])
	set(context, "zone", load["time_zone"])
	set(context, "ntp-server", load["ntp_server"])
}

func set(context *cli.Context, key, val string) {
	if context.String(key) == "" {
		if err := context.Set(key, val); err != nil {
			logrus.Fatalln(err)
		}
	}
}
