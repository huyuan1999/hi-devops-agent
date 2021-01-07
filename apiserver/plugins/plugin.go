package plugins

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"gopkg.in/ini.v1"
	"io"
	"log"
	"os/exec"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PLUGIN",
	MagicCookieValue: "HI-DEVOPS",
}

type Plugin struct {
	LoggerOutput io.Writer
	LoggerLevel  hclog.Level
}

func (p *Plugin) logger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: p.LoggerOutput,
		Level:  p.LoggerLevel,
	})
}

func (p *Plugin) config(name string) (map[string]string, error) {
	pl := make(map[string]string)
	//cfg, err := ini.Load(config.CfgPath)
	cfg, err := ini.Load("./test.ini")
	if err != nil {
		return nil, err
	}


	pl["name"] = cfg.Section(fmt.Sprintf("plugin:%s", name)).Key("name").String()
	pl["path"] = cfg.Section(fmt.Sprintf("plugin:%s", name)).Key("path").String()

	if pl["name"] == "" || pl["path"] == "" {
		return nil, errors.New("plugin config name or path is empty")
	}
	return pl, nil
}

func (p *Plugin) loadWithContext(pluginMap map[string]plugin.Plugin, pluginPath string) *plugin.Client {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(pluginPath),
		Logger:          p.logger(),
	})
	return client
}

func (p *Plugin) Load(tag string) error {
	pl, err := p.config(tag)
	if err != nil {
		return err
	}

	var pluginMap = map[string]plugin.Plugin{
		pl["name"]: &RegisterPlugin{},
	}

	client := p.loadWithContext(pluginMap, pl["path"])

	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	raw, err := rpcClient.Dispense(pl["name"])
	if err != nil {
		return err
	}

	register_plugins := raw.(RegisterPlugins)
	log.Println("register_plugins.Register(): ", register_plugins.Register())
	return nil
}
