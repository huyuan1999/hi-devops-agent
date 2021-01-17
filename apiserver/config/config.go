package config

var (
	ConfigFile string
	OnlyId     string
	NtpServer  string
	TimeZone   string
)

var (
	InitializationFile = ".init.json"
	CfgServer          string
	CfgWorkdir         string
	CfgAddress         string // grpc 简单地址, 格式 ip:port
	CfgPublicIP        string // register 向服务器汇报数据时标识 agent 节点的地址
	CfgDaemon          bool
	CfgPidFile         string
	CfgLogFile         string
	CfgSetTime         bool
	CfgSetZone         bool
)

var (
	ServerRequestRegister = "/api/v1/install/agent/register/"
	ServerRequestCASign   = "/api/v1/tls/ca/sign/"
)

var (
	CertFile = "tls/server.pem"
	CertCsr  = "tls/server.csr"
	KeyFile  = "tls/server.key"
	CAPem    = "tls/ca.pem"
)
