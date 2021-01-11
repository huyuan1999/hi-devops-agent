package linux

import (
	"github.com/huyuan1999/hi-devops-agent/cmdb/linux/command"
	"github.com/huyuan1999/hi-devops-agent/cmdb/models"
	"github.com/huyuan1999/hi-devops-agent/cmdb/utils"
)

type MainBoard struct {
	models.MainBoard
	dmidecode *command.Dmidecode
}

func NewMainBoard() (*MainBoard, error) {
	mainBoard := &MainBoard{}
	dmidecode, err := command.NewDmidecode()
	if err != nil {
		return nil, err
	}
	mainBoard.dmidecode = dmidecode
	utils.Call(mainBoard)
	return mainBoard, nil
}

func (m *MainBoard) GetSerialNumber() {
	m.SerialNumber = m.dmidecode.SystemSerialNumber()
}

func (m *MainBoard) GetUUID() {
	m.UUID = m.dmidecode.SystemUuid()
}

func (m *MainBoard) GetManufacturer() {
	m.Manufacturer = m.dmidecode.SystemManufacturer()
}

func (m *MainBoard) GetProductName() {
	m.ProductName = m.dmidecode.SystemProductName()
}
