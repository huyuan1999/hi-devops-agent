package services

import (
	"context"
	"github.com/huyuan1999/hi-devops-agent/cmdb/linux"
)

type ClientService struct {
}

func (c *ClientService) GetSystem(ctx context.Context, request *RequestClient) (*System, error) {
	system, err := linux.NewSystem()
	if err != nil {
		return nil, err
	}
	return &System{
		HostName: system.HostName,
		Release:  system.Release,
		Kernel:   system.Kernel,
	}, nil
}

func (c *ClientService) GetCPU(ctx context.Context, request *RequestClient) (*CPU, error) {
	cpu, err := linux.NewCPU()
	if err != nil {
		return nil, err
	}
	return &CPU{
		Number:    uint32(cpu.Number),
		Core:      uint32(cpu.Core),
		Sibling:   uint32(cpu.Sibling),
		Processor: uint32(cpu.Processor),
		ModelName: cpu.ModelName,
	}, nil
}

func (c *ClientService) GetMemory(ctx context.Context, request *RequestClient) (*Memory, error) {
	mem, err := linux.NewMemory()
	if err != nil {
		return nil, err
	}
	return &Memory{
		Total:    uint64(mem.Total),
		Type:     mem.Type,
		Number:   uint32(mem.Number),
		Slot:     uint32(mem.Slot),
		MaxSize:  mem.MaxSize,
		FreeSlot: uint32(mem.FreeSlot),
	}, nil
}

func (c *ClientService) GetMainBoard(ctx context.Context, request *RequestClient) (*MainBoard, error) {
	board, err := linux.NewMainBoard()
	if err != nil {
		return nil, err
	}
	return &MainBoard{
		SerialNumber: board.SerialNumber,
		Uuid:         board.UUID,
		Manufacturer: board.Manufacturer,
		ProductName:  board.ProductName,
	}, nil
}

func (c *ClientService) GetNIC(ctx context.Context, request *RequestClient) (*NICMany, error) {
	nicSlice, err := linux.NewNIC()
	if err != nil {
		return nil, err
	}

	var n []*NICMany_NICOne

	for _, nic := range nicSlice {
		nicOne := &NICMany_NICOne{
			Name:    nic.Name,
			Mac:     nic.Mac,
			Address: nic.Address,
		}
		n = append(n, nicOne)
	}

	return &NICMany{
		Nic: n,
	}, nil
}

func (c *ClientService) GetDisk(ctx context.Context, request *RequestClient) (*DiskMany, error) {
	diskSlice, err := linux.NewDisk()
	if err != nil {
		return nil, err
	}

	var d []*DiskManyDiskOne

	for _, disk := range diskSlice {
		diskOne := &DiskManyDiskOne{
			Name:         disk.Name,
			SerialNumber: disk.SerialNumber,
			ProductName:  disk.ProductName,
			Size:         uint64(disk.Size),
			FormFactor:   disk.FormFactor,
		}
		d = append(d, diskOne)
	}
	return &DiskMany{
		Disk: d,
	}, nil
}
