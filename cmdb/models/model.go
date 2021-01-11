package models


type System struct {
	HostName string `json:"host_name"`
	Release  string `json:"release"`
	Kernel   string `json:"kernel"`
}

type CPU struct {
	Number    uint   `json:"number"`
	Core      uint   `json:"core"`
	Sibling   uint   `json:"sibling"`
	Processor uint   `json:"processor"` // 服务器总线程数 = Number * Siblings
	ModelName string `json:"model_name"`
}

type Memory struct {
	Total    uint   `json:"total"`
	Type     string `json:"type"`
	Number   uint   `json:"number"`
	Slot     uint   `json:"slot"`
	MaxSize  string `json:"max_size"`
	FreeSlot uint   `json:"free_slot"`
}

type MainBoard struct {
	SerialNumber string `json:"serial_number"`
	UUID         string `json:"uuid"`
	Manufacturer string `json:"manufacturer"`
	ProductName  string `json:"product_name"`
}

type NIC struct {
	Name    string   `json:"name"`
	Mac     string   `json:"mac"`
	Address []string `json:"address"`
}

type Disk struct {
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
	ProductName  string `json:"product_name"`
	Size         uint   `json:"size"`
	FormFactor   string `json:"form_factor"`
}
