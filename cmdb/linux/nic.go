package linux

import (
	"github.com/huyuan1999/hi-devops-agent/cmdb/models"
	"github.com/huyuan1999/hi-devops-agent/cmdb/utils"
	"net"
)

type NIC []struct {
	models.NIC
}

func NewNIC() (NIC, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	nic := make(NIC, len(ifaces))
	utils.Call(nic)
	return nic, nil
}

func (n NIC) GetName() {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for index, iface := range ifaces {
		n[index].Name = iface.Name
	}
}

func (n NIC) GetMac() {
	ifaces, err := net.Interfaces()
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
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for index, iface := range ifaces {
		n[index].Mac = iface.HardwareAddr.String()
	}
}
