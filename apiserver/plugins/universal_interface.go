package plugins

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

type RegisterPlugins interface {
	Register() error
}

type RegisterRPC struct{ client *rpc.Client }

func (r *RegisterRPC) Register() error {
	var e error
	err := r.client.Call("Plugin.Register", new(interface{}), &e)
	if err != nil {
		return err
	}
	if e != nil {
		return e
	}
	return nil
}

type RegisterPlugin struct {
	Impl RegisterPlugins
}

func (p *RegisterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RegisterPlugin{Impl: p.Impl}, nil
}

func (RegisterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RegisterRPC{client: c}, nil
}
