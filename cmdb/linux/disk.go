package linux

import (
	"errors"
	"github.com/huyuan1999/hi-devops-agent/cmdb/linux/command"
	"github.com/huyuan1999/hi-devops-agent/cmdb/models"
	"github.com/huyuan1999/hi-devops-agent/cmdb/utils"
	"github.com/shirou/gopsutil/v3/disk"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"syscall"
	"unsafe"
)

type Disk []struct {
	models.Disk
	diskNameArr []string
	smart       *command.Smart
}

const diskstats = "/proc/diskstats"

func isDisk(array []string) (diskArray []string) {
	for _, base := range array {
		for _, dev := range array {
			if dev == base {
				continue
			}
			if strings.HasPrefix(dev, base) {
				diskArray = append(diskArray, base)
				goto loop
			}
		}
	loop:
	}
	return
}

func NewDisk() (Disk, error) {
	text, err := ioutil.ReadFile(diskstats)
	if err != nil {
		return nil, err
	}
	compile, err := regexp.Compile("\\s+8\\s+\\d+\\s+\\w+")
	if err != nil {
		return nil, err
	}

	diskArr := compile.FindAllString(string(text), -1)
	var nameArr []string
	for _, dev := range diskArr {
		split := strings.Split(utils.DeleteExtraSpace(utils.Trim(dev)), " ")
		if len(split) >= 3 {
			name := path.Join("/dev/", split[2])
			nameArr = append(nameArr, name)
		}
	}
	diskNameArr := isDisk(nameArr)
	diskNumber := len(diskNameArr)
	if diskNumber < 1 {
		return nil, errors.New("获取硬盘信息错误: 硬盘数量小于 1")
	}
	smart, _ := command.NewSmart()
	d := make(Disk, diskNumber)
	d[0].diskNameArr = diskNameArr
	d[0].smart = smart
	utils.Call(d)
	return d, nil
}

func (d Disk) GetName() {
	for index, dev := range d[0].diskNameArr {
		d[index].Name = dev
	}
}

func (d Disk) ioctl(device string, event uintptr, value uintptr) error {
	if fd, err := unix.Open(device, os.O_RDONLY, 0660); err != nil {
		return err
	} else {
		_, _, ErrOn := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), event, value)
		if unix.ErrnoName(ErrOn) != "" {
			return errors.New(ErrOn.Error())
		}
	}
	return nil
}

func (d Disk) serialNumberWithContext(device string) string {
	// 使用 smartctl 命令
	if d[0].smart != nil {
		if serial := d[0].smart.SerialNumber(device); serial != "" {
			return serial
		}
	}

	// 读取 /run/udev/data/ 信息
	serial, err := disk.SerialNumber(device)
	if err == nil && serial != "" {
		s := strings.Split(serial, "_")
		return s[len(s)-1]
	}

	// 使用 ioctl 函数
	var hd unix.HDDriveID
	if err := d.ioctl(device, unix.HDIO_GET_IDENTITY, uintptr(unsafe.Pointer(&hd))); err == nil {
		var sn []byte
		for _, char := range hd.Serial_no {
			sn = append(sn, char)
		}
		if string(sn) != "" {
			return string(sn)
		}
	}
	return ""
}

func (d Disk) GetSerialNumber() {
	for index, dev := range d[0].diskNameArr {
		d[index].SerialNumber = d.serialNumberWithContext(dev)
	}
}

func (d Disk) productWithContext(device string) string {
	if d[0].smart != nil {
		if manufacturer := d[0].smart.DeviceModel(device); manufacturer != "" {
			return manufacturer
		}
	}

	serial, err := disk.SerialNumber(device)
	if err == nil && serial != "" {
		if s := strings.Split(serial, "_"); len(s)-2 >= 0 {
			model := s[0 : len(s)-2]
			return strings.Join(model, "_")
		}
	}

	var hd unix.HDDriveID
	if err := d.ioctl(device, unix.HDIO_GET_IDENTITY, uintptr(unsafe.Pointer(&hd))); err == nil {
		var model []byte
		for _, char := range hd.Model {
			model = append(model, char)
		}
		if string(model) != "" {
			return string(model)
		}
	}
	return ""
}

func (d Disk) GetManufacturer() {
	for index, dev := range d[0].diskNameArr {
		d[index].ProductName = d.productWithContext(dev)
	}
}

func (d Disk) GetSize() {
	var size uint64
	for index, dev := range d[0].diskNameArr {
		fd, err := unix.Open(dev, os.O_RDONLY, 0660)
		if err != nil {
			continue
		}
		_, _, ErrOn := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), unix.BLKGETSIZE64, uintptr(unsafe.Pointer(&size)))
		if unix.ErrnoName(ErrOn) != "" {
			continue
		}
		d[index].Size = uint(size >> 30)
	}
}

func (d Disk) GetFormFactor() {
	for index, dev := range d[0].diskNameArr {
		d[index].FormFactor = d[0].smart.FormFactor(dev)
	}
}
