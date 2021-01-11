package command

import (
	"github.com/huyuan1999/hi-devops-agent/cmdb/utils"
)

type Dmidecode struct {
}

func NewDmidecode() (*Dmidecode, error) {
	if err := utils.Which("dmidecode"); err != nil {
		return nil, err
	} else {
		return &Dmidecode{}, nil
	}
}

func (d *Dmidecode) Info() string {
	result := utils.Shell("dmidecode")
	return utils.Trim(result.Stdout)
}

func (d *Dmidecode) SystemManufacturer() string {
	result := utils.Shell("dmidecode", "-s", "system-manufacturer")
	return utils.Trim(result.Stdout)
}

func (d *Dmidecode) SystemProductName() string {
	result := utils.Shell("dmidecode", "-s", "system-product-name")
	return utils.Trim(result.Stdout)
}

func (d *Dmidecode) SystemSerialNumber() string {
	result := utils.Shell("dmidecode", "-s", "system-serial-number")
	return utils.Trim(result.Stdout)
}

func (d *Dmidecode) SystemUuid() string {
	result := utils.Shell("dmidecode", "-s", "system-uuid")
	return utils.Trim(result.Stdout)
}
