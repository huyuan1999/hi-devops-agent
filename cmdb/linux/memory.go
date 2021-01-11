package linux

import (
	"github.com/huyuan1999/hi-devops-agent/cmdb/linux/command"
	"github.com/huyuan1999/hi-devops-agent/cmdb/models"
	"github.com/huyuan1999/hi-devops-agent/cmdb/utils"
	"github.com/shirou/gopsutil/mem"
	"regexp"
	"strings"
)

type Memory struct {
	models.Memory
	dmidecodeInfo string
}

func NewMemory() (*Memory, error) {
	memory := &Memory{}
	dmidecode, err := command.NewDmidecode()
	if err != nil {
		return nil, err
	}
	memory.dmidecodeInfo = dmidecode.Info()
	utils.Call(memory)
	return memory, nil
}

func (m *Memory) GetTotal() {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return
	}
	m.Total = uint(memory.Total)
}

func (m *Memory) GetType() {
	matched, err := utils.LoopMatchString(m.dmidecodeInfo, []string{"(?s)(?U)Memory\\s+Device\\n+.*Memory\\s+Device\\n+", "Type:\\s+.*"})
	if err != nil {
		return
	}
	dataArr := strings.Split(matched, ":")
	if len(dataArr) == 2 {
		m.Type = utils.Trim(dataArr[1])
	}
}

func (m *Memory) GetNumber() {
	compile, err := regexp.Compile("Memory\\s+Device\\s+Mapped\\s+Address")
	if err != nil {
		return
	}
	matched := compile.FindAllString(m.dmidecodeInfo, -1)
	m.Number = uint(len(matched))
}

func (m *Memory) GetSlot() {
	compile, err := regexp.Compile("(?s)Memory\\s+Device\n+.*Configured")
	if err != nil {
		return
	}
	memInfo := compile.FindString(m.dmidecodeInfo)
	re, err := regexp.Compile("Size:\\s.*")
	if err != nil {
		return
	}
	slot := re.FindAllString(memInfo, -1)
	m.Slot = uint(len(slot))
}

func (m *Memory) GetMaxSize() {
	matched, err := utils.LoopMatchString(m.dmidecodeInfo, []string{"Maximum\\sCapacity:.*", "\\d+\\s+GB"})
	if err != nil {
		return
	}
	m.MaxSize = matched
}

func (m *Memory) GetFreeSlot() {
	compile, err := regexp.Compile("Size:\\sNo\\sModule\\sInstalled")
	if err != nil {
		return
	}

	slot := compile.FindAllString(m.dmidecodeInfo, -1)
	m.FreeSlot = uint(len(slot))
}
