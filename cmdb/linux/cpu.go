package linux

import (
	"github.com/huyuan1999/hi-devops-agent/cmdb/models"
	"github.com/huyuan1999/hi-devops-agent/cmdb/utils"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type CPU struct {
	models.CPU
	document string
}

const cpuinfo = "/proc/cpuinfo"

func NewCPU() (*CPU, error) {
	cpu := &CPU{}
	text, err := ioutil.ReadFile(cpuinfo)
	if err != nil {
		return nil, err
	}
	cpu.document = string(text)
	utils.Call(cpu)
	return cpu, nil
}

func (c *CPU) GetNumber() {
	compile, err := regexp.Compile("physical\\s+id.*")
	if err != nil {
		return
	}
	physical := compile.FindAllString(c.document, -1)
	c.Number = uint(len(utils.RemoveDuplicate(physical)))
}

func (c *CPU) GetCore() {
	matched, err := utils.LoopMatchString(c.document, []string{"cpu\\s+cores.*", "\\d+"})
	if err != nil {
		return
	}

	if core, err := strconv.Atoi(matched); err == nil {
		c.Core = uint(core)
	}
}

func (c *CPU) GetSibling() {
	matched, err := utils.LoopMatchString(c.document, []string{"siblings.*", "\\d+"})
	if err != nil {
		return
	}

	if sibling, err := strconv.Atoi(matched); err == nil {
		c.Sibling = uint(sibling)
	}
}

func (c *CPU) GetProcessor() {
	c.GetNumber()
	c.GetSibling()
	c.Processor = c.Number * c.Sibling
}

func (c *CPU) GetModelName() {
	compile, err := regexp.Compile("model\\s+name.*")
	if err != nil {
		return
	}
	modelName := compile.FindString(c.document)

	modelArr := strings.Split(modelName, ":")
	if len(modelArr) == 2 {
		c.ModelName = utils.Trim(modelArr[1])
	}
}
