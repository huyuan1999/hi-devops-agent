package linux

import (
	"github.com/huyuan1999/hi-devops-agent/cmdb/models"
	"github.com/huyuan1999/hi-devops-agent/cmdb/utils"
	"io/ioutil"
	"net"
)

const virtualPath = "/sys/devices/virtual/net/"

type NIC []struct {
	models.NIC
}

// 过滤虚拟网卡
func virtualNet() ([]net.Interface, error) {
	var virtualNameSlice []string
	var newInterface []net.Interface

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	virtual, err := ioutil.ReadDir(virtualPath)
	if err != nil {
		return nil, err
	}

	for _, item := range virtual {
		virtualNameSlice = append(virtualNameSlice, item.Name())
	}
	for _, iface := range ifaces {
		if !utils.IsContain(virtualNameSlice, iface.Name) {
			newInterface = append(newInterface, iface)
		}
	}
	return newInterface, nil
}

func NewNIC() (NIC, error) {
	ifaces, err := virtualNet()
	if err != nil {
		return nil, err
	}
	nic := make(NIC, len(ifaces))
	utils.Call(nic)
	return nic, nil
}

func (n NIC) GetName() {
	ifaces, err := virtualNet()
	if err != nil {
		return
	}
	for index, iface := range ifaces {
		n[index].Name = iface.Name
	}
}

func (n NIC) GetMac() {
	ifaces, err := virtualNet()
	if err != nil {
		return
	}
	for index, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			if ip = ip.To4(); ip == nil {
				continue
			}
			n[index].Address = append(n[index].Address, ip.String())
		}
	}
}

func (n NIC) GetAddress() {
	ifaces, err := virtualNet()
	if err != nil {
		return
	}
	for index, iface := range ifaces {
		n[index].Mac = iface.HardwareAddr.String()
	}
}
