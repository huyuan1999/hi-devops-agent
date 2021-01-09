package main

import (
	"errors"
	"github.com/huyuan1999/hi-devops-agent/apiserver/config"
	"github.com/huyuan1999/hi-devops-agent/apiserver/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/ini.v1"
	"plugin"
	"strings"
)

type Register struct {
	pluginName string
	pluginPath string
}

func NewRegister() *Register {
	return &Register{}
}

func (r *Register) config(cfg *ini.File, name string) error {
	r.pluginName = cfg.Section(name).Key("name").String()
	r.pluginPath = cfg.Section(name).Key("path").String()

	if r.pluginName == "" || r.pluginPath == "" {
		return errors.New("plugin config name or path is empty")
	}
	return nil
}

func (r *Register) loadWithContext() {
	cfg, err := ini.Load(config.CfgPath)
	if err != nil {
		log.Errorf("load ini config error: %s", err.Error())
		return
	}

	var pluginsName []string
	var pluginsPath []string

	for _, section := range cfg.Sections() {
		pluginMap := make(map[string]string)

		if !strings.HasPrefix(section.Name(), "plugin") {
			continue
		}

		if err := r.config(cfg, r.pluginName); err != nil {
			log.Errorf("parsing plugin %s configuration error: %s", r.pluginName, err.Error())
			continue
		}

		if utils.IsContain(pluginsName, r.pluginName) {
			log.Warningf("plugin %s has been loaded", r.pluginName)
			continue
		}

		if utils.IsContain(pluginsPath, r.pluginPath) {
			log.Warningf("the same plugin %s is loaded repeatedly", r.pluginPath)
			continue
		}

		p, err := plugin.Open(r.pluginPath)
		if err != nil {
			log.Errorf("open plugin %s error: %s", r.pluginName, err.Error())
			continue
		}

		// 每一个插件都必须在 main 包中实现 Register 方法,
		symbol, err := p.Lookup("Register")
		if err != nil {
			log.Errorf("lookup plugin %s error: %s", r.pluginName, err.Error())
			continue
		}

		// Register 方法只能接收一个参数, 并参数名为 server 类型为 *grpc.Server
		if err := symbol.(func(server *grpc.Server) error)(config.GRPCServer); err == nil {
			pluginMap["name"] = r.pluginName
			pluginMap["path"] = r.pluginPath
			pluginsName = append(pluginsName, r.pluginName)
			pluginsPath = append(pluginsPath, r.pluginPath)
			config.PluginInfo = append(config.PluginInfo, pluginMap)
		} else {
			log.Errorf("call plugin %s error: %s", r.pluginName, err.Error())
		}
	}
}

func (r *Register) Load() {
	r.loadWithContext()
}
